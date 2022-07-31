FROM node:16.14.2 as node-build
COPY web /app/
WORKDIR /app
ARG branchName="DEFAULT_BRANCH_NAME"
ARG commitId="DEFAULT_COMMIT_ID"
ENV branchName=${branchName}
ENV commitId=${commitId}
RUN sed -i "s|UnknownBranchName|${branchName}|g" ./src/environments/environment.prod.ts
RUN sed -i "s|UnknownCommitId|${commitId}|g" ./src/environments/environment.prod.ts
RUN rm -fr /app/dist/*
RUN ls -al /app/dist
RUN npm install
RUN npm run build-prod
RUN ls -al /app/dist

FROM golang:1.18.4-alpine as go-build
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN rm -fr static
COPY --from=node-build /app/dist/web /usr/src/app/static
RUN ls -al
RUN ls -al static
RUN go version
RUN go build -o app
RUN ls -al

FROM golang:1.18.4-alpine
WORKDIR /usr/src/app
COPY --from=go-build /usr/src/app/app .
RUN ls -al
EXPOSE 8080
CMD ["/usr/src/app/app"]
