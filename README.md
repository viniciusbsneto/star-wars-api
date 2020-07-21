

<h1 align="center">
    <img alt="Star Wars API" title="#Star Wars API" src="https://upload.wikimedia.org/wikipedia/commons/6/6c/Star_Wars_Logo.svg" />
</h1>

<h1 align="center">
   🚀 <a href="#"> STAR WARS API </a> 🌌
</h1>

<h3 align="center">
    A Star Wars API. Find out more about the planets!
</h3>

<p align="center">
  <img alt="GitHub language count" src="https://img.shields.io/github/languages/count/viniciusbsneto/star-wars-api?color=yellow">

  <img alt="Repository size" src="https://img.shields.io/github/repo-size/viniciusbsneto/star-wars-api">
  
  <a href="https://github.com/viniciusbsneto/star-wars-api/commits/master">
    <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/viniciusbsneto/star-wars-api">
  </a>
    
   <img alt="License" src="https://img.shields.io/badge/license-MIT-yellow">
   <a href="https://github.com/viniciusbsneto/star-wars-api/stargazers">
    <img alt="Stargazers" src="https://img.shields.io/github/stars/viniciusbsneto/star-wars-api?style=social">
  </a>

  <a href="https://github.com/viniciusbsneto">
    <img alt="made by viniciusbsneto" src="https://img.shields.io/badge/-viniciusbsneto-yellow">
  </a>
</p>


<h4 align="center"> 
	 Status: Core finished 🚧 Refactoring... 🚧
</h4>

<p align="center">
 <a href="#about">About</a> •
 <a href="#features">Features</a> •
 <a href="#pending-to-do">Pending (to do)</a> •
 <a href="#how-it-works">How it works</a> • 
 <a href="#tech-stack">Tech Stack</a> •  
 <a href="#author">Author</a> • 
 <a href="#user-content-license">License</a>
</p>


## :speech_balloon: About

Star Wars API - is a technical challenge proposed by an awesome company for a junior back-end position that I'm applying to.
It's a REST API for Star Wars planets.

---

## :bulb: Features

- [x] Create a planet
  - [x] Get the number of film appearances of a planet by cosuming an external API
  - [SWAPI](https://swapi.dev/)
- [x] Update a planet
- [x] Delete a planet
- [x] Get all planets
- [x] Get a planet by ID
- [x] Get a planet by name

---

## :hammer_and_wrench: Pending (to do) :hourglass_flowing_sand:
1. [ ] **Error handling**
	- I've spent most of the time on learning Go language and how to make a REST API with it.
	- I've also spent some time on studying MongoDB driver for Go to setup database storage for the API.
2. [ ] **Project structure**
	- I don't know what project structure I should use with a Go REST API. I have considered MVC or DDD.
	- But now I'm currently studying a Go Packages structure approach as it is said to be the standard strucutre among Go developers.
3. [ ] **Tests**
	- I'm currently learning how to write tests in Go. I should be pushing some tests very soon.
4. [ ] **Database container**
	- The API uses a MongoDB Atlas cluster database (MongoDB Atlas cluster connection URL). I should write a Dockerfile to create a container for MongoDB image so others don't need to use MongoDB Atlas cluster.
5. [ ] **Application image**
	- Finally I'm considering making an image of the whole application and its dependencies with Docker so it's easier for others to run it.

## ⚙ How it works

For now this project is comprised of these files:
1. god.mod - For dependencies
2. god.sum - For checksums
3. .env - For environment variables (MongoDB Atlas URI for now)
4. main.go (src folder) - All API code resides in this file (I plan on refactoring)

### :pushpin: Pre-requisites

Before you begin, you will need to have the following tools installed on your machine:
[Git] (https://git-scm.com), [Go] (https://golang.org/).
In addition, it is good to have an editor to work with the code like [VSCode] (https://code.visualstudio.com/)

#### Rodando o Backend (servidor)

```bash

# Clone this repository
$ git clone git@github.com:viniciusbsneto/star-wars-api.git

# Access the project folder cmd/terminal
$ cd star-wars-api

# go to the src folder
cd src

# run the application
$ go run main.go

# The server will start at port: 3333 - go to http://localhost:3333

```

<p align="center">
  <a href="https://github.com/viniciusbsneto/star-wars-api/blob/master/Insomnia_API_Star_Wars.json" target="_blank"><img src="https://insomnia.rest/images/run.svg" alt="Run in Insomnia"></a>
</p>

#### Endpoints
Do NOT type curly braces

POST /planets - Create a planet
  - To create a planet send a POST request to /planets endpoint with JSON body (ID is generated automatically by MongoDB):
  ```
  {
    "name": "Alderaan",
    "climate" [
      "temperate"
    ],
    "terrain": [
      "grasslands",
      "mountain"
    ]
  }
  ```
PUT /planets/{id} - Update a planet
  - To update a planet send a PUT request to /planets endpoint with new JSON body and passing an ID in the route:
   ```
  {
    "name": "NEW Alderaan",
    "climate" [
      "temperate",
      "a new climate",
    ],
    "terrain": [
      "grasslands",
      "mountain"
    ],
    "films": 2
  }
  ```
DEL /planets/{id} - Delete a planet
  - To delete a planet send a DEL request to /planets passing an ID in the route:
  ```
  Ex.: DEL http://localhost:3333/planets/5f157effe08b5ad3ffa598e6
  ```
GET /planets - Get all planets
  - To get all planets send a GET request to /planets endpoint
  ```
  Ex.: GET http://localhost:3333/planets
  ```
GET /planets/{id} - Get a planet by ID
  - To get a planet by ID send a GET request to /planets passing an ID in the route:
  ```
  Ex.: GET http://localhost:3333/planets/5f157effe08b5ad3ffa598e6
  ```
GET /search?name={planetName} - Get a planet by Name
  - To get a planet by name send a GET request to /search endpoint passing the query to a planet name:
  ```
  Ex.: GET http://localhost:3333/search?name=Alderaan
  ```

---

## :toolbox: Tech Stack

The following tools were used in the construction of the project:

#### [](https://github.com/viniciusbsneto/star-wars-api#rest-api)**REST API**

-   **[Go](https://golang.org/)**
-   **[Gorilla Mux](https://github.com/gorilla/mux)**
-   **[MongoDB](https://www.mongodb.com/)**
-   **[GoDotEnv](https://github.com/joho/godotenv)**

> See the file  [go.mod](https://github.com/viniciusbsneto/star-wars-api/blob/master/go.mod)

#### [](https://github.com/viniciusbsneto/star-wars-api#utilit%C3%A1rios)**Utilitários**

-   Editor:  **[Visual Studio Code](https://code.visualstudio.com/)**
-   Markdown:  **[StackEdit](https://stackedit.io/)**,  **[Markdown Emoji](https://gist.github.com/rxaviers/7360908)**
-   API Test:  **[Insomnia](https://insomnia.rest/)**

---

## :handshake: How to contribute

1. Fork the project.
2. Create a new branch with your changes: `git checkout -b my-feature`
3. Save your changes and create a commit message telling you what you did: `git commit -m" feature: My new feature "`
4. Submit your changes: `git push origin my-feature`
> If you have any questions check this [guide on how to contribute](./CONTRIBUTING.md)

---

## :technologist: Author

 <img style="border-radius: 50%;" src="https://avatars1.githubusercontent.com/u/17788722?v=4" width="100px;" alt="Vinícius Neto"/> 
 <br />

[![Linkedin Badge](https://img.shields.io/badge/-Vinícius%20Neto-blue?style=flat-square&logo=Linkedin&logoColor=white&link=https://www.linkedin.com/in/vinicius-neto/)](https://www.linkedin.com/in/vinicius-neto/) 
[![Gmail Badge](https://img.shields.io/badge/-viniciusbsneto@gmail.com-c14438?style=flat-square&logo=Gmail&logoColor=white&link=mailto:viniciusbsneto@gmail.com)](mailto:viniciusbsneto@gmail.com)

---

## :memo: License

This project is under the license [MIT](./LICENSE).

Made with love by Vinícius Neto 👋🏽 [Get in Touch!](Https://www.linkedin.com/in/vinicius-neto/)
