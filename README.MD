## About The Project

Create Spotify playlists from Radio Helsinki programs.


## Getting Started


### Prerequisites

Spotify developer documentation [https://developer.spotify.com/documentation/](https://developer.spotify.com/documentation/) 

Install `docker`

### Running

1. With docker
   ```
   docker build -t . radiohelsinki-to-spotify
   docker run -p 8080:8080 -e REDIRECT_URL=http://localhost:8080/callback -e SPOTIFY_ID=<get from Spotify dashboard> -e SPOTIFY_SECRET=<get from Spotify dashboard> radiohelsinki-to-spotify
   ```

2. Without docker
   ```
   export SPOTIFY_ID=<get from Spotify dashboard>
   export SPOTIFY_SECRET=<get from Spotify dashboard>
   export REDIRECT_URL=http://localhost:8080/callback
   go mod download
   go run main.go
   ```  

## Usage

Visit [http://localhost:8080](http://localhost:8080)

Login with your Spotify account and press the "Create playlist" button.



## Roadmap

* Support more RH programs