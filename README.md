# Book Management software
A simple book management software.


# Product Requirements

**1. Add and manage books:** The software should allow users to add books to the system and manage them. It should capture basic information about each book, such as title, author, published date, edition, description, genre, and any other relevant details.

**2. Create and manage collections of books:** Users should be able to create collections of books. This feature enables organizing books into different categories or groups based on specific criteria, allowing for better organization and easy retrieval.

**3. List books and collections:** The software should provide the ability to list all the books and collections present in the system. This feature allows users to get an overview of the available books and collections.

**4. Filter book lists:** Users should be able to apply filters to book lists based on various criteria. These filters can include author, genre, or a range of publication dates. Filtering helps users narrow down their search and find specific books based on their preferences.


   
# APIs

### 1. Add a Book

- **Endpoint**: `/api/v1/books`
- **Method**: `POST`
- **Request Payload**:
```json
{
  "title": "The Book Title",
  "author": "Author Name",
  "published_date": "2023-05-01",
  "edition": "1st Edition",
  "description": "Book description",
  "genre": "Fiction"
}
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
