###########################################################
# Building stage
###########################################################
FROM golang:1.13 as build

RUN apt-get update && apt-get install -y asciidoctor && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace
COPY . .

RUN asciidoctor docs/release-notes.adoc

ENV GOOS=linux
ENV GOARCH=386

RUN go build
RUN for CMD in `ls cmd`; do go build -o ./build/$CMD ./cmd/$CMD ; done


###########################################################
# Running stage
###########################################################
FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app/

COPY --from=build /workspace/config-api .
COPY --from=build /workspace/build/* ./
COPY --from=build /workspace/database/migrations ./database/migrations
COPY --from=build /workspace/config/*.toml ./config/
COPY --from=build /workspace/docs/* ./docs/
COPY --from=build /workspace/mta_version .

CMD ["sh", "-c", "./dbSync && ./config-api"]
