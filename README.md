# search-app

## Overview
`search-app` is a web application that allows users to perform text searches and view the results. The application is built using Vue.js for the frontend and Go for the backend. This application is designed to showcase the full-text search capabilities of Cloud Spanner.

## Project Setup

### Install Dependencies
To install the necessary dependencies, run:
```sh
npm install
```

### Compiles and Hot-Reloads for Development
To start the development server with hot-reloading, run:

```sh
npm run serve
```

### Compiles and Minifies for Production

To build the project for production, run:

```sh
npm run build
```

## Backend Setup

### Install Go Dependencies

Ensure you have Go installed. Then, install the necessary Go dependencies by running:

```sh
go mod tidy
```

### Import sample data to Cloud Spanner

Create a Spanner instance, database and table.

```sql
CREATE TABLE Restaurants (
  id STRING(MAX) NOT NULL,
  dateAdded TIMESTAMP OPTIONS (
    allow_commit_timestamp = true
  ),
  dateUpdated TIMESTAMP OPTIONS (
    allow_commit_timestamp = true
  ),
  address STRING(MAX),
  categories STRING(MAX),
  primaryCategories STRING(MAX),
  city STRING(MAX),
  country STRING(MAX),
  keys STRING(MAX),
  latitude FLOAT64,
  longitude FLOAT64,
  name STRING(MAX),
  postalCode STRING(MAX),
  province STRING(MAX),
  sourceURLs STRING(MAX),
  websites STRING(MAX),
  name_token TOKENLIST AS (tokenize_fulltext(name)) HIDDEN,
  categories_token TOKENLIST AS (tokenize_substring(categories)) HIDDEN,
  city_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(city)) HIDDEN,
) PRIMARY KEY(id);;
```

Create the index for full text search
```sql
CREATE SEARCH INDEX RestaurantsIndex ON Restaurants(name_token, categories_token);
```

Use the sample data [Fast Food Restaurants Across America](https://data.world/datafiniti/fast-food-restaurants-across-america) and import it into Cloud Spanner to demonstrate its full-text search capabilities.

```sh
go run main.go -import -file=Datafiniti_Fast_Food_Restaurants_Jun19.csv
```

### Run the Backend Server

To start the backend server, run:

```sh
go run main.go
```

## Frontend Setup

### Using the Application

Performing a Search

1. Open the application in your web browser.
2. Enter your search query in the search bar.
3. Click the "Search" button to perform the search.

### Viewing Search Results
The search results will be displayed below the search bar. Each result includes details such as country, city, name, address, websites, and categories.

## Debugging

### Print Executed Query
The application prints the executed SQL query to the console for debugging purposes. This helps in verifying the correctness of the query.

### Print Search Results
The application also prints each search result to the console. This helps in verifying the correctness of the search results.

## Directory Structure

project-root/
├── main.go
├── public/
│   ├── index.html
├── src/
│   ├── assets/
│   │   └── tailwind.css
│   ├── App.vue
│   └── main.js
├── babel.config.js
├── postcss.config.js
├── tailwind.config.js
├── package.json
├── README.md
└── vue.config.js

## Additional Information

### Tailwind CSS
The project uses Tailwind CSS for styling. Ensure that the Tailwind CSS configuration is correctly set up in tailwind.config.js and postcss.config.js.

### Vue Configuration
The Vue configuration is defined in vue.config.js, which sets the development server port and other settings.