<p align="center">
  <a href="https://www.codechefvit.com" target="_blank">
    <img src="https://i.ibb.co/4J9LXxS/cclogo.png" width=160 title="CodeChef-VIT" alt="Codechef-VIT">
  </a>
</p>

<h2 align="center"> CookOff 10.0 Backend </h2>

---

> CookOff is CodeChef VIT's flagship competitive coding event that tests the coding skills of developers.  
> This backend powers both the **admin** and **participant portals** for CookOff 10.0, handling users, questions, test cases, submissions, leaderboards, and timers.  
> Designed for scalability and reliability, it ensures seamless competition management and smooth participation.



## Tech Stack

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [SQLC](https://github.com/sqlc-dev/sqlc)

---

## Features

### User Management

- `POST /signup` – User signup
- `POST /login` – Login
- `POST /logout` – Logout
- `POST /refreshToken` – Refresh tokens
- `GET /dashboard` – Load participant dashboard

### Question Management

- `GET /admin/questions` – Get all questions
- `GET /admin/questions/:id` – Get question by ID
- `POST /admin/questions` – Create a question
- `PUT /admin/questions/:id` – Update a question
- `DELETE /admin/questions/:id` – Delete a question
- `POST /admin/questions/:id/bounty/activate` – Activate bounty
- `POST /admin/questions/:id/bounty/deactivate` – Deactivate bounty

### Testcase Management

- `GET /testcase/:id` – Get a testcase by ID
- `GET /testcase` – Get all testcases
- `GET /question/:id/testcases` – Get all testcases for a question
- `GET /question/:id/testcases/public` – Get only public testcases
- `POST /testcase` – Create testcase (Admin only)
- `PUT /testcase/:id` – Update testcase (Admin only)
- `DELETE /testcase/:id` – Delete testcase (Admin only)

### Submission Management

- `POST /submit` – Submit code
- `POST /runcode` – Run code against hidden testcases
- `POST /runcustom` – Run custom input
- `GET /result/:submission_id` – Get submission result

### Leaderboard

- `GET /leaderboard` – Fetch leaderboard

### Timer

- `POST /admin/setTime` – Set round time
- `POST /admin/updateTime` – Update round time
- `GET /admin/startRound` – Start round
- `GET /admin/resetRound` – Reset round
- `GET /getTime` – Get remaining time

### Admin User Controls

- `GET /admin/users` – Get all users
- `POST /admin/users/:id/ban` – Ban user
- `POST /admin/users/:id/unban` – Unban user
- `POST /admin/users/:id/upgrade` – Upgrade user to next round
- `GET /admin/users/:id/submissions` – Get all submissions by a user

---

## Getting Started

### Installation

1. Fork the repo
2. Clone it locally
   ```sh
   git clone https://github.com/<GITHUB_USERNAME>/cookoff-10.0-be.git
   cd cookoff-10.0-be
   ```

### Prerequisites

- Configure environment variables (`.env` file – refer `.env.example`)
- Configure Makefile

### Running Locally

1. Start containers
   ```sh
   docker compose up --build -d
   ```

2. Install SQLC
   ```sh
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

3. Generate SQLC schema and queries
   ```sh
   make generate
   ```

4. Apply migrations
   ```sh
   make up 
   ```

---

## Related Projects

- **Cookoff Admin Portal:** https://github.com/CodeChefVIT/cookoff-admin-10.0
- **Cookoff Portal:** https://github.com/CodeChefVIT/cookoff-portal-10.0
---


[![UI](https://img.shields.io/badge/User%20Interface-Link%20to%20UI-orange?style=flat-square&logo=appveyor)](https://cookoff24.codechefvit.com/)

---

## Contributors

<table>
  <tr align="center">
    <td>
      <p align="center">
        <a href="https://github.com/Soham-Maha">
          <img src="https://avatars.githubusercontent.com/u/155614230?v=4" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Soham Mahapatra</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/equestrian2296">
          <img src="https://avatars.githubusercontent.com/equestrian2296" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Atharva Sharma</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/BrainNotFoundException">
          <img src="https://avatars.githubusercontent.com/BrainNotFoundException" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Lavnish Jhunjunwala</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/aayushk231">
          <img src="https://avatars.githubusercontent.com/aayushk231" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Aayush Kushwaha</p>
    </td>
  </tr>
  <tr align="center">
    <td>
      <p align="center">
        <a href="https://github.com/Advik-Gupta">
          <img src="https://avatars.githubusercontent.com/Advik-Gupta" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Advik Gupta</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/RustyDev24">
          <img src="https://avatars.githubusercontent.com/RustyDev24" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Vedant Matanhelia</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/Shrish2006">
          <img src="https://avatars.githubusercontent.com/Shrish2006" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Shrish</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/ASHUTOSH-SWAIN-GIT">
          <img src="https://avatars.githubusercontent.com/ASHUTOSH-SWAIN-GIT" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Ashutosh Swain</p>
    </td>
  </tr>
  <tr align="center">
    <td>
      <p align="center">
        <a href="https://github.com/upayanmazumder">
          <img src="https://avatars.githubusercontent.com/upayanmazumder" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Upayan Mazumder</p>
    </td>
    <td>
      <p align="center">
        <a href="https://github.com/abhitrueprogrammer">
          <img src="https://avatars.githubusercontent.com/abhitrueprogrammer" width="170" height="170" style="border:2px solid grey; border-radius:50%;">
        </a>
      </p>
      <p style="font-size:17px; font-weight:600;">Abhinav Pant</p>
    </td>
  </tr>
</table>

---

## License

[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)

---

<p align="center">
  Made with ❤️ by <a href="https://www.codechefvit.com" target="_blank">CodeChef-VIT</a>
</p>
