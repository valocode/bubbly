# Fetch Repo Licenses
## Install go-licenses 
```go get github.com/google/go-licenses```

## Run go-licenses
```go-licenses csv ./...```

# Fetch Repo CVEs

## Install Snyk

```npm install -g snyk ```

## Authenticate Snyk Account
```snyk auth```

## Run Snyk Test
```snyk test --json-file-output=./snyk.json```