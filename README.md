# 🎬 Filmigo REST API

A powerful REST API for movie data built with Go and the Filmigo library. Get movie information from IMDB, OMDB, and JustWatch through simple HTTP endpoints.

## 🚀 Live Demo

**API Base URL:** `https://your-app-name.onrender.com`

Try these endpoints:
- **Health Check:** [/api/v1/health](https://your-app-name.onrender.com/api/v1/health)
- **Get Movie:** [/api/v1/movies/tt0111161](https://your-app-name.onrender.com/api/v1/movies/tt0111161)
- **Search Movies:** [/api/v1/movies/search?q=inception](https://your-app-name.onrender.com/api/v1/movies/search?q=inception)

## 📋 API Endpoints

| Method | Endpoint | Description | Example |
|--------|----------|-------------|---------|
| GET | `/api/v1/health` | API health check | - |
| GET | `/api/v1/movies/{id}` | Get movie by IMDB ID | `tt0111161` |
| GET | `/api/v1/movies/search?q={query}` | Search movies | `?q=inception` |
| GET | `/api/v1/omdb/{id}` | Get OMDB movie data | `tt0111161` |

## 🌟 Features

- ✅ Get detailed movie information by IMDB ID
- ✅ Search movies by title (requires OMDB API key)
- ✅ Access OMDB data with enhanced details
- ✅ Fast response times with built-in caching
- ✅ CORS enabled for web applications
- ✅ Comprehensive error handling
- ✅ Professional JSON responses


### Get Movie Details
