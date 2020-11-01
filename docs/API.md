---
title: Article similarity v1.0.0
language_tabs: []
language_clients: []
toc_footers: []
includes: []
search: false
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="article-similarity">Article similarity v1.0.0</h1>

> Scroll down for example requests and responses.

Server to store articles and search similar articles.

Base URLs:

* <a href="/">/</a>

<h1 id="article-similarity-default">Default</h1>

## post__articles

`POST /articles`

*Add an article.*

> Body parameter

```json
{
  "content": "string"
}
```

<h3 id="post__articles-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» content|body|string|true|Article content|

> Example responses

> 201 Response

<h3 id="post__articles-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Article added.|[Article](#schemaarticle)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid arguments|[Error](#schemaerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## get__articles

`GET /articles`

*Get unique articles.*

> Example responses

> 200 Response

<h3 id="get__articles-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK.|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[Error](#schemaerror)|

<h3 id="get__articles-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» articles|[[Article](#schemaarticle)]|true|none|none|
|»» id|[ArticleId](#schemaarticleid)(int64)|true|none|Article id|
|»» content|string|true|none|Article content|
|»» duplicate_article_ids|[integer]|true|none|Duplicated articles|

<aside class="success">
This operation does not require authentication
</aside>

## get__articles_{id}

`GET /articles/{id}`

*Get article by id.*

<h3 id="get__articles_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|integer(int64)|true|Article id|

> Example responses

> 200 Response

<h3 id="get__articles_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[Article](#schemaarticle)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid arguments|[Error](#schemaerror)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Article not found.|[Error](#schemaerror)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## get__duplicate_groups

`GET /duplicate_groups`

*Get duplicate groups ids.*

> Example responses

> 200 Response

<h3 id="get__duplicate_groups-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK.|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[Error](#schemaerror)|

<h3 id="get__duplicate_groups-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» duplicate_groups|[array]|true|none|none|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_Error">Error</h2>
<!-- backwards compatibility -->
<a id="schemaerror"></a>
<a id="schema_Error"></a>
<a id="tocSerror"></a>
<a id="tocserror"></a>

```json
{
  "code": 0,
  "message": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|integer(int64)|false|none|Error code for machine parsing|
|message|string|true|none|Human-readable error message|

<h2 id="tocS_ArticleId">ArticleId</h2>
<!-- backwards compatibility -->
<a id="schemaarticleid"></a>
<a id="schema_ArticleId"></a>
<a id="tocSarticleid"></a>
<a id="tocsarticleid"></a>

```json
0

```

Article id

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|integer(int64)|false|none|Article id|

<h2 id="tocS_Article">Article</h2>
<!-- backwards compatibility -->
<a id="schemaarticle"></a>
<a id="schema_Article"></a>
<a id="tocSarticle"></a>
<a id="tocsarticle"></a>

```json
{
  "id": 0,
  "content": "string",
  "duplicate_article_ids": [
    0
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|[ArticleId](#schemaarticleid)|true|none|Article id|
|content|string|true|none|Article content|
|duplicate_article_ids|[integer]|true|none|Duplicated articles|

