# Scalability

"Article similarity" system consists of two elements:
- web server
- `mongodb` database

## Vertical scaling

System can be scaled vertically without restrictions. It is needed to add processing power and HTTP requests 
would be handled faster.

## Horizontal scaling

Web server is stateless and designed to be horizontally scalable. `mongodb` supports sharding. So, just write as many
app and database containers as needed and system could process more requests.
