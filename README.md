# Book Management Software
The Book Management Package is a Go package designed to facilitate the management of books and collections of books. It provides a set of APIs for adding books, creating collections, listing books and collections, and filtering book lists based on various criteria. The package also includes a well-defined database schema for storing book and collection information.



# Features

**1. Add and manage books:** The software should allow users to add books to the system and manage them. It should capture basic information about each book, such as title, author, published date, edition, description, genre, and any other relevant details.

**2. Create and manage collections of books:** Users should be able to create collections of books. This feature enables organizing books into different categories or groups based on specific criteria, allowing for better organization and easy retrieval.

**3. List books and collections:** The software should provide the ability to list all the books and collections present in the system. This feature allows users to get an overview of the available books and collections.

**4. Filter book lists:** Users should be able to apply filters to book lists based on various criteria. These filters can include author, genre, or a range of publication dates. Filtering helps users narrow down their search and find specific books based on their preferences.

# Usages
## 1. Adding a book to the system

```curl
curl -X POST -H "Content-Type: application/json" -d '{
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965-08-01",
  "edition": "1st Edition",
  "description": "Paul Muad'Dib leads the Fremen on a conquest of revenge",
  "genre": "Science Fiction"
}' /api/books

```
- **Example Response**:
```json
{
  "book_id": "1234",
  "staus": "success", 
  "code": "200"
}
```

- **Example Error Response**:
```json
{
  "staus": "error", 
  "code": "400"
}
```
## 2. Creating a collection
```curl
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "Dune",
  "description": "The collected sayings of Muad'Dib (by the Princess Irulan)."
}' /api/collections

```

- **Example Response**:

```json
{
  "collection_id": "5678",
  "status": "success", 
  "code" : "200"
}

```

- **Example Error Response**:
```json
{
  "staus": "error", 
  "code": "400"
}
```
# APIs

## 1. Add a Book

- **Endpoint**: `/api/v1/books`
- **Description**: This endpoint allows you to add a book to the book management system. You need to provide the necessary information about the book, such as the title, author, published date, edition, description, and genre. After successfully adding the book, it will return the book ID and a status code.
- **Method**: `POST`
- **Request Payload**: 
```json
{
  "title": "Dune",
  "author": "Frank Herbert",
  "published_date": "1965-08-01",
  "edition": "1st Edition",
  "description": "Paul Muad'Dib leads the Fremen on a conquest of revenge",
  "genre": "Science Fiction"
}
```
- **Response**:
```json
{
  "book_id": "1234",
  "staus": "success", 
  "code" : "200"
}

```

## 2. Create a Collection

- **Endpoint**: `/api/v1/collections`
- **Description**: This endpoint enables you to create a collection of books. You can specify the name and description of the collection. Upon successful creation, it will return the collection ID and a status code.
- **Method**: `POST`
- **Request Payload**:
```json
{
  "name": "Dune",
  "description": "The collected sayings of Muad'Dib (by the Princess Irulan)."
}
```
- **Response**:

```json
{
  "collection_id": "5678",
  "status": "success", 
  "code" : "200"
}

```
## 3. List Books

- **Endpoint**: `/api/v1/books`
- **Description**: This endpoint allows you to retrieve a list of all the books in the system. It returns an array of book objects, each containing information such as the book ID, title, author, published date, edition, description, and genre. Use this endpoint to get an overview of all available books.
- **Method**: `GET`
- **Response**:
```json
[
  {
    "book_id": "1234",
    "title": "Dune",
    "author": "Frank Herbert",
    "published_date": "1965-08-01",
    "edition": "1st Edition",
    "description": "Paul Muad'Dib leads the Fremen on a conquest of revenge",
    "genre": "Science Fiction"
  },
  ...
]
```

## 4. List Collections
- **Endpoint**: `/api/v1/collections`
- **Description**: This endpoint allows you to retrieve a list of collections from the system.
- **Method**: `GET`
- **Response**:

```json
[
  {
    "collection_id": "5678",
    "name": "Dune",
    "description": "The collected sayings of Muad'Dib (by the Princess Irulan)."
  },
  ...
]
```

## 5. Filter Books
- **Endpoint**: `/api/books/filter`
- **Description**: This endpoint allows you to filter book lists by author, genre, or a range of publication dates.
- **Method**: `GET`
- **Query Parameters**:
  - `author`: Filter books by author name.
  - `genre`: Filter books by genre.
  - `from_date`: Filter books published from a specific date.
  - `to_date`: Filter books published until a specific date.
- **Example**:
  ```bash
  curl -X GET '/api/books/filter?title=Dune&genre=Science%20Fiction&from_date=1960-01-01&to_date=1970-12-31'

  ```
- **Response**:
```json
[
  {
    "book_id": "1234",
    "title": "Dune",
    "author": "Frank Herbert",
    "published_date": "1965-08-01",
    "edition": "1st Edition",
    "description": "Paul Muad'Dib leads the Fremen on a conquest of revenge",
    "genre": "Science Fiction"
  },
  ...
]

```

# Database Schema

### Books Table

| Column Name     | Data Type    | Description                                    |
| --------------- | -------------| ---------------------------------------------- |
| book_id         | Primary Key  | Unique identifier for the book                  |
| title           |  String      | Title of the book                              |
| author          |   String     | Author of the book                             |
| published_date  |   Date       | Publication date of the book                    |
| edition         |    Int       | Edition of the book                             |
| description     |    String    | Description of the book                         |
| genre           |    String    | Genre of the book                               |
| ...             |              | (Additional columns as needed for relevant details) |

### Collections Table

| Column Name     | Data Type    | Description                                    |
| --------------- | -------------| ---------------------------------------------- |
| collection_id   | Primary Key  | Unique identifier for the collection            |
| name            |  String      | Name of the collection                          |
| description     |  String      | Description of the collection                   |

### CollectionBooks Table (Many-to-Many Relationship)

| Column Name     | Data Type    | Description                                    |
| --------------- | -------------| ---------------------------------------------- |
| collection_id   | Foreign Key  | References the collection_id in Collections table|
| book_id         | Foreign Key  | References the book_id in Books table           |
