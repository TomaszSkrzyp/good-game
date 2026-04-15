
**WIP**
# 🏀 **Good Game**

**Good Game** is a **Go-based web application** designed for NBA fans who want to find the **best games to watch** without having the results spoiled. Whether you're looking for a close finish or a historic performance, the app calculates a proprietary **Quality Score** for every game. It leverages a **high-concurrency backend** and a **PostgreSQL database** to aggregate real-time data from the NBA ecosystem.

With **automated data pipelines**, **Good Game** syncs with external APIs every 15 minutes to track live scores and schedules. Using a **spoiler-free UI philosophy**, the application hides scores and winners by default, allowing users to choose what to watch based on excitement levels rather than outcomes.

---

## 📌 **Features**

- **Proprietary Quality Algorithm:** Automatically rates games from 0-100 based on game dynamic.
- **Spoiler-Free Browsing:** A "Hide Scores" mode is baked into the core architecture to ensure you never see a result before the game.
- **Automated Sync Worker:** A dedicated background service that handles high-frequency updates without affecting API performance.
- **Clutch Metrics:** Highlights games with one-possession finishes and massive fourth-quarter comebacks.
- **Responsive Frontend:** Built with **SolidJS** for fine-grained reactivity and a lightning-fast user experience.
- **Containerized Environment:** Fully orchestrated with **Docker** for consistent deployment across environments.

---

## 🚀 **Installation**

### **Prerequisites**
- **Go 1.24+**
- **Node.js 22+**
- **Docker** and **Docker Compose**
- **PostgreSQL 15+**
---
### **1. Clone the Repository**
Clone the repository to your local machine:

```bash
git clone [https://github.com/tomaszSkrzyp/good-game.git](https://github.com/tomaszSkrzyp/good-game.git)
cd good-game
```
### **2. Configure Environment Variables**
Create a .env file in the project root with the following:

```bash
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=goodgame
DATABASE_URL="host=db user=postgres password=your_password dbname=goodgame port=5432 sslmode=disable"
JWT_SECRET=your_super_secret_key
PORT=8080
```

### **3. Launch via Docker Compose**
The easiest way to run the full stack (API, Frontend, Database, and Worker) is using Docker:

```bash
docker-compose up --build
```

### **4. Run a Manual Season Sync (Optional)**
If you want to populate the database with historical data for the 2026 season immediately:

```bash
docker exec -it gg-backend ./main --update-all
```
### **How the Game Score is Calculated**
A 0-100 Points badge will be displayed for each game. The rules for scoring are as follows:
- **Final Margin**: 
  - **Margin <3 pts: +45 pts**
  - **Margin <7 pts: +30 pts**
  - **Margin <12 pts: +15 pts**
- **Clutch Factor**: Tight scores entering the 4th quarter or a final margin <3: +20 pts
- **Huge Swing**: Overcoming a 15+ point deficit in the 4th quarter: +25 pts
- **Star Duel**: Elite scorers (35+ pts) on both teams going head-to-head: +20 pts
- **Big Game**: Individual performance of 30+ pts across 3+ statistical categories: +15 pts
- **Overtime**: Any game reaching extra periods: +15 pts
- **Style Points**: High-scoring "Shootouts" or defensive "Gritty" battles: +10 to +15 pts

### **Badge Definitions**
Badges provide a quick visual cue of the game's personality without revealing the winner:
- **🌪️ IsHugeSwing**: A massive comeback or lead change occurred in the second half.
- **🎯 IsClutch**: The game was decided in the final possessions.
- **⚔️ IsStarDuel**: Two or more superstars had legendary performances against each other.
- **🚀 IsShootou**t: An ultra-high-scoring affair (Total > 235 pts).
- **🛡️ IsGritty**: A defensive, physical battle (Total < 200 pts).
- **🔥 IsBigGame** : A player had a historic individual stat line (e.g., a 50-pt game)

### **User Ratings & Community**
While the algorithm provides an objective "Watchability" score, Good Game also captures the human element:
- **Aggregated Sentiment**: The app displays an Average User Rating next to the algorithmic score.
- **Spoiler-Free Voting**: Users can rate games they have watched
- **Rating Synchronization**: User reactions are stored with unique constraints (per User/Game ID) to ensure data integrity and prevent double-voting.


### **🛠 Technologies Used**
- **Go (Golang)** – High-performance backend used for the API and the background data-fetching worker.
- **SolidJS**  – A declarative, efficient UI library with fine-grained reactivity for the frontend.
- **PostgreSQL** & GORM – Advanced relational database paired with the GORM library for complex data relationships.
- **Docker & Compose** – Used for containerization and seamless orchestration of the microservice architecture.
- **JWT (JSON Web Tokens)** – Handles stateless authentication for user profiles and game ratings.
- **Tailwind CSS** – Custom utility-first styling for a polished, responsive, and modern dark-themed UI.
- **ESPN API Integration** – Provides the raw data feed for schedules, box scores, and play-by-play updates.

### **📜 License**
This project is licensed under the MIT License, allowing for free use, modification, and distribution.

### **🔗 Contributing**
I welcome contributions to enhance the scoring algorithm or the UI!
Fork the repository to your GitHub account.
Create a new feature branch:
```bash
git checkout -b feature/amazing-feature
```
Make your changes and commit them with a descriptive message.
Push to the branch and open a Pull Request.

### **📧 Contact**
For questions, bug reports, or suggestions, please reach out via the  page. We value your feedback!
