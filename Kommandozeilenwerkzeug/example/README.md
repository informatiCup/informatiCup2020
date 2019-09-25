# With Docker

## Build

```
docker build --tag icup2020_example .
```

## Run

```
docker run --publish 50123:50123 --interactive --tty --rm icup2020_example
```

# Without Docker

## Setup

```
pip install --user -r requirements.txt
```

## Run

```
./example.py
```
