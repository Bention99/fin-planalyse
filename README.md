# FinPlanalyse

FinPlanalyse is a personal finance analysis and budgeting application
built with Go and PostgreSQL.\
It helps users track income and expenses, categorize transactions, and
analyze financial behavior in a structured and scalable way.

------------------------------------------------------------------------

## 💡 Motivation

Tracking your budget down to each individual transaction can be a lot of work, especially if your bank doesn’t provide meaningful analysis or flexible categorization options.

I built Fin-Planalyse because I wanted deeper insight into my real spending behavior. Instead of relying on static banking dashboards, I wanted a tool that allows me to:

-   Track every transaction in a structured way
-   Categorize income and expenses flexibly
-   Analyze where my money actually goes
-   Plan how to optimize and get the best out of my budget

FinPlanalyse gives you full control over your financial data. It combines detailed transaction tracking with structured analysis, helping turn raw bank statements into actionable financial insight.

------------------------------------------------------------------------

## 🚀 Features

-   📂 Upload and extract transactions from bank statements (PDF)
-   🤖 AI-powered transaction extraction using OpenAI
-   💰 Income & Expense categorization
-   🏷 System and user-defined categories
-   📊 Basic financial overview and analysis
-   🔐 User-based data separation
-   🗄 PostgreSQL database
-   🔄 Goose migrations for schema management

------------------------------------------------------------------------

## 🏗 Tech Stack

-   **Backend:** Go (net/http)
-   **Database:** PostgreSQL
-   **Migrations:** Goose
-   **SQL Codegen:** sqlc
-   **AI Integration:** OpenAI API
-   **Frontend:** HTML / CSS

------------------------------------------------------------------------

## ⚙️ Quick Start

### 1. Clone the repository

``` bash
git clone https://github.com/Bention99/fin-planalyse
cd finplanalyse
```

### 2. Configure environment variables

Create a `.env` file:

    DATABASE_URL=postgres://user:password@localhost:5432/finplanalyse?sslmode=disable
    OPENAI_API_KEY=your_api_key_here

### 3. Run database migrations

``` bash
goose postgres "$DATABASE_URL" up
```

### 4. Run the application

``` bash
go run .
```

App runs at:

    http://localhost:8080

------------------------------------------------------------------------

## 🎮 Usage

-   Create an account
-   Manage categories (system or custom)
-   Add transactions manually
-   Upload bank statements (PDF) for automatic extraction
-   Analyze your income and expenses

------------------------------------------------------------------------

## 🤖 AI Transaction Extraction

FinPlanalyse integrates with OpenAI to:

-   Extract transactions from PDF statements
-   Classify transactions into categories
-   Structure output into JSON

The uploaded file is always stored in:

    /uploads/statement.pdf

Only one file is stored at a time.

------------------------------------------------------------------------

## 🛣 Roadmap

-   Authentication & session management
-   Financial dashboards
-   Monthly budget targets
-   Savings tracking
-   Deployment setup (Docker)

------------------------------------------------------------------------

## 📄 License

This project is currently for educational and personal use.

------------------------------------------------------------------------

## 🪄 Contributing

If you have ideas, improvement suggestions, or would like to contribute, feel free to open an issue or submit a pull request. Any feedback is welcome, whether it’s about features, code structure, performance, or overall design.