package submissions

import (
    "context"
    "log"
    "math/big"

    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

const SubmissionDoneStatus = "DONE"
type ResultResp struct {
    ID             string     `json:"id"`
    QuestionID     string     `json:"question_id"`
    Passed         int        `json:"passed"`
    Failed         int        `json:"failed"`
    Runtime        float64    `json:"runtime"`
    Memory         float64    `json:"memory"`
    SubmissionTime string     `json:"submission_time"`
    Description    string     `json:"description"`
    Testcases      []TCResult `json:"testcases"`
}
type TCResult struct {
    ID          string  `json:"id"`
    Runtime     float64 `json:"runtime"`
    Memory      float64 `json:"memory"`
    Status      string  `json:"status"`
    Description string  `json:"description"`
}

func CheckStatus(ctx context.Context, subID uuid.UUID) (bool, error) {
    status, err := utils.Queries.GetSubmissionStatusByID(ctx, subID)
    if err != nil {
        return false, err
    }
    if status == nil {
        return false, nil
    }
    return *status == SubmissionDoneStatus, nil
}

func GetSubResult(ctx context.Context, subID uuid.UUID) (ResultResp, error) {
    submission, err := utils.Queries.GetSubmissionByID(ctx, subID)
    if err != nil {
        return ResultResp{}, err
    }

    subResults, err := utils.Queries.GetSubmissionResultsBySubmissionIDQuery(ctx, subID)
    if err != nil {
        return ResultResp{}, err
    }

    var desc string
    if submission.Description != nil {
        desc = *submission.Description
    }

    subRuntime, _ := submission.Runtime.Float64Value()
    subMemory, _ := submission.Memory.Float64Value()

    resp := ResultResp{
        ID:             submission.ID.String(),
        QuestionID:     submission.QuestionID.String(),
        Passed:         int(submission.TestcasesPassed.Int32),
        Failed:         int(submission.TestcasesFailed.Int32),
        Runtime:        subRuntime.Float64,
        Memory:         subMemory.Float64,
        SubmissionTime: submission.SubmissionTime.Time.String(),
        Description:    desc,
        Testcases:      make([]TCResult, len(subResults)),
    }

    for i, result := range subResults {
        runtime, _ := result.Runtime.Float64Value()
        memory, _ := result.Memory.Float64Value()

        tcID := ""
        if result.TestcaseID.Valid {
            tcID = result.TestcaseID.UUID.String()
        }

        resultDesc := ""
        if result.Description != nil {
            resultDesc = *result.Description
        }

        resp.Testcases[i] = TCResult{
            ID:          tcID,
            Runtime:     runtime.Float64,
            Memory:      memory.Float64,
            Status:      result.Status,
            Description: resultDesc,
        }
    }

    return resp, nil
}
func UpdateSubmission(ctx context.Context, subID uuid.UUID) error {
    status := SubmissionDoneStatus

    data, err := utils.Queries.GetStatsForFinalSubEntryBySubmissionID(ctx, subID)
    if err != nil {
        log.Println("Error fetching submission results:", err)
        return err
    }

    var totalRuntime float64
    var totalMemory int64
    var passed, failed int

    for _, v := range data {
        runtime, _ := v.Runtime.Float64Value()
        totalRuntime += runtime.Float64
        totalMemory += v.Memory.Int.Int64()
        if v.Status == "success" {
            passed++
        } else {
            failed++
        }
    }
    err = utils.Queries.UpdateSubmissionByID(ctx, db.UpdateSubmissionByIDParams{
        ID:              subID,
        Runtime:         pgtype.Numeric{Int: big.NewInt(int64(totalRuntime)), Valid: true},
        Memory:          pgtype.Numeric{Int: big.NewInt(totalMemory), Valid: true},
        Status:          &status,
        TestcasesPassed: pgtype.Int4{Int32: int32(passed), Valid: true},
        TestcasesFailed: pgtype.Int4{Int32: int32(failed), Valid: true},
    })
    if err != nil {
        log.Println("Error updating submission:", err)
        return err
    }


    if err := utils.Queries.UpdateUserScoreBySubmissionID(ctx, subID); err != nil {
        log.Println("Error updating user score:", err)
        return err
    }

    log.Printf("Submission updated successfully: %v\n", subID)
    return nil
}
