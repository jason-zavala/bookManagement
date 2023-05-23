# API usage examples

## 1. Adding a book to the system


``` bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965",
  "edition": "First Edition",
  "description": "A science fiction epic set in a distant future where interstellar politics, religion, and ecology intertwine.",
  "genre": "Science Fiction"
}
' http://localhost:8080/api/v1/books

```

``` bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "title": "Harry Potter and the Sorcerers Stone",
  "author": "J.K. Rowling",
  "published_date": "1997",
  "edition": "First Edition",
  "description": "The magical journey begins as Harry Potter discovers he is a wizard.",
  "genre": "Fantasy"
}
' http://localhost:8080/api/v1/books

```
``` bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "title": "The Hitchhikers Guide to the Galaxy",
  "author": "Douglas Adams",
  "published_date": "1979",
  "edition": "First Edition",
  "description": "A delightful jaunt through the cosmos where a clueless human, a perpetually depressed robot, and a book with all the answers join forces in a hilarious quest to explore the universe, dodge Vogons, and ponder the meaning of life",
  "genre": "Science Fiction"
}
' http://localhost:8080/api/v1/books

```

``` bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "title": "The Hunger Games",
  "author": "Suzanne Collins",
  "published_date": "2008",
  "edition": "First Edition",
  "description": "In a dystopian future, Katniss Everdeen fights for survival in a televised battle royale against other oppressed teens.",
  "genre": "Science Fiction"
}
' http://localhost:8080/api/v1/books

```

## 2. Creating a Collection

```bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "name": "Dune",
  "description": "The collected sayings of MuadDib (by the Princess Irulan)."
}
' http://localhost:8080/api/v1/collections
```
```bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "name": "Harry Potter",
  "description": " A tale of a wizard prodigy who attends a magical school, battles dark forces, and narrowly escapes death on a yearly basis, all while his friends discover their knack for causing mischief."
}
' http://localhost:8080/api/v1/collections
```
```bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "name": "The Hitchhikers Guide to the Galaxy",
  "description": "An ordinary human is whisked away on a bewildering intergalactic journey, encountering eccentric aliens, a depressed robot, and discovering the ultimate answer to life, the universe, and everything "
}
' http://localhost:8080/api/v1/collections
```
```bash
curl -X POST -H "Content-Type: application/json" -d '
{
  "name": "The Hunger Games",
  "description": "In a dystopian future, a courageous young girl from a poverty-stricken district volunteers for a televised fight to the death against other oppressed teens."
}
' http://localhost:8080/api/v1/collections
```
