
# Firefly Home Assignment

Please read the instructions to run the project locally


## Intro

Since the original assignment includes a large number of essays, I created "mini" lists of essays and word bank that is a subset of the original lists. Running the progran with the short list is of course optional. When the program starts, it will ask you whether you would like to run the short or full version
## Run Locally

Clone the project

```bash
  git clone https://github.com/ronihm/firefly.git
```

Go to the project directory

```bash
  cd firefly
```

**Make sure you have words.txt in your project's root directory**. You can achieve this by one of 2 options:
Either manually put the file in the directory, or pull it from git lfs (see below)

Now build the docker image
```bash
  docker build -t fireflyapp .
```

Run the docker with **interactive mode**

```bash
  docker run -i fireflyapp
```


## Pulling words.txt

Since words.txt is bigger than the limit file size allowed by github, I used git lfs. To pull words.txt from git lfs follow the following steps:

1. Download git lfs from the [official website](https://git-lfs.com/)
2. From your project directory, run:


```bash
  git lfs pull
```
## What I Did

First I created the bank of allowed words. I chose a trie over a map for better space efficency.
Then I created 2 worker pools - one for fetching the essays and the other for counting the words on each essay.
In order to avoid being blocked by the server, I used the following mechanisms:
1. Rate limiting of the requests (using a threads-safe rate limiter)
2. Random user-agents
3. If we are still getting blocked, there is a retry mechanism that waits an exponentially growing amount of time between tries.

To count the words I used a threads-safe map (sync.Map). I wrapped the functionality in a WordCounter class for clarity.
After all workers have finished processing the essays, I'm using a heap to get the 10 most popular words