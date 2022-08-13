FROM node:16.14.2 as node-build
COPY web /app/
WORKDIR /app
ARG branchName="DEFAULT_BRANCH_NAME"
ARG commitId="DEFAULT_COMMIT_ID"
ENV branchName=${branchName}
ENV commitId=${commitId}
RUN ls -al
RUN sed -i "s|UnknownBranchName|${branchName}|g" ./src/environments/environment.prod.ts
RUN sed -i "s|UnknownCommitId|${commitId}|g" ./src/environments/environment.prod.ts
RUN rm -fr /app/dist/*
RUN npm install
RUN npm run build-prod

FROM golang:1.18.4-alpine as go-build
ARG branchName="DEFAULT_BRANCH_NAME"
ARG commitId="DEFAULT_COMMIT_ID"
ENV branchName=${branchName}
ENV commitId=${commitId}
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN rm -fr static
COPY --from=node-build /app/dist/web /usr/src/app/static
RUN echo "$(cat version.txt)-${commitId}" > version.txt
RUN cat version.txt
RUN go build -o app

FROM alpine:3.16
WORKDIR /usr/src/app
COPY --from=go-build /usr/src/app/app .
EXPOSE 8080
CMD ["/usr/src/app/app"]
