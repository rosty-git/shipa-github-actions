# shipa-github-actions

Shipa github actions backend uses in create app and deploy app actions

## Build docker image

    docker build -t vmanilo/shipa-action:0.1.6 .


### Run action

     docker run --env SHIPA_HOST=<host> --env SHIPA_TOKEN=<token> vmanilo/shipa-action:0.1.6.3