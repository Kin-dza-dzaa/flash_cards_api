## Overview

This project was built for learning purposes. The API lets you to manage collections of words with translations provided by [google.translate.com.](google.translate.com)   
For authorization/autontification opendID connect is used. Connection with [google.translate.com](google.translate.com) established through HTTP 2.0. 

It's only the backend of the whole application. The application itself can be found at: [https://github.com/Kin-dza-dzaa/flash_cards](https://github.com/Kin-dza-dzaa/flash_cards) 

See the table of documentation:

<table><tbody><tr><td>End point</td><td>method</td><td>Request body</td><td>Response</td><td>Description</td></tr><tr><td>/v1/words</td><td>GET</td><td>empty</td><td><pre><code class="language-javascript">{ 
   “path”: “/v1/words”, 
   “status”: 200, 
   "user_words”: {
      “collection_name”: {
          [translation data]  
      }
   }
} </code></pre></td><td>Gets all user words from each collection.</td></tr><tr><td>/v1/words</td><td>POST</td><td><pre><code class="language-javascript">{
   “word”: string,
   “collection_name”: string,
   “last_repeat”: string(time),
   “time_diff”: int64 timespan 
}</code></pre></td><td>200 OK</td><td>Add the word to the specified collection.&nbsp;</td></tr><tr><td>/v1/words</td><td>DELETE</td><td><pre><code class="language-javascript">{
   “word”: string,
   “collection_name”: string
}</code></pre></td><td>200 OK</td><td>Deletes word from collection.</td></tr><tr><td>/v1/words</td><td>PUT</td><td><pre><code class="language-javascript">{
   “word”: string,
   “collection_name”: string,
   “last_repeat”: string(time),
   “time_diff”: int64 timespan 
}</code></pre></td><td>200 OK</td><td>Updates learn interval of the provided word.</td></tr></tbody></table>

## Usage

**Run app:**

```plaintext
make run
```

**Run tests:**

**In order for tests to work you will need an available docker API on port 2375 with disabled tls.**

```plaintext
make test
```

**Run tests with cover:**

```plaintext
make cover
```
