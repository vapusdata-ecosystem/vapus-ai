## **Service Overview**

The `AIUtility` service offers Remote Procedure Calls (RPCs):

1. **GenerateEmbedding**  
   Generates vector embeddings for the provided text.

2. **SensitivityAnalyzer**  
   Analyzes and optionally acts on sensitive entities found in the text.

---

## **Authentication**

All requests to the `AIUtility` service require an authorization token (JWT). Include the token in the `Authorization` header as a **Bearer Token**:


## **How to Use**

### 1. **GenerateEmbedding**

**RPC Name:** `GenerateEmbedding`  
**Description:** Generates vector embeddings for the given text.

#### **Request Format**

Send a `GenerateEmbeddingRequest` message containing a list of strings (`text`) for which embeddings are required.

Example:
```json
{
  "text": ["example sentence 1", "example sentence 2"]
}
```

### 2. **SensitivityAnalyzer**

**RPC Name:** `SensitivityAnalyzer`  
**Description:** Analyzes and acts on the given text by identifying sensitive entities based on the specified action and post-detection action.

#### **Request Format**

Send a `SensitivityAnalyzerRequest` message containing:
- `text`: A list of strings to analyze.
- `entities`: A list of entity keywords that need to be detected.
- `action`: The action to perform on the text (e.g., `ANALYZE` or `ACT`).
- `postDetectAction`: The post-detection action to perform on the identified entities (e.g., `REDACT`, `FAKE`, `EMPTY`).

Example:
```json
{
  "text": ["This is a sensitive example text."],
  "entities": ["sensitive", "example"],
  "action": "ANALYZE",
  "postDetectAction": "REDACT"
}
```
