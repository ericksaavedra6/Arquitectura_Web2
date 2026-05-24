# Banco de Pruebas: gRPC (localhost:50051)
### Nota para Postman: Asegúrate de seleccionar el método gRPC, importar tu archivo product.proto en la pestaña Service Definition y elegir el método correspondiente en el menú desplegable.

## Método: Create
### Cuerpo del mensaje (JSON):

{
  "nombre": "Reloj de Titanio Saavedra",
  "descripcion": "Modelo Elegance de 38mm con caja sostenible",
  "precio": 450.00
}

## Método: Get (Existente)
### Cuerpo del mensaje (JSON):

{
  "id": 1
}
## Método: Get (Manejo de Errores - Inexistente)
### Cuerpo del mensaje (JSON):

{
  "id": 999
}
(Debe retornar un código de estado 4 NOT_FOUND con tu mensaje personalizado).
## Método: Update
### Cuerpo del mensaje (JSON):

{
  "id": 1,
  "nombre": "Reloj de Titanio Saavedra Elegance v2",
  "descripcion": "Caja de titanio de 38mm y calibre automático modificado",
  "precio": 480.00
}
## Método: List
### Cuerpo del mensaje (JSON):

{}
## Método: Delete
### Cuerpo del mensaje (JSON):

{
  "id": 1
}
# Banco de Pruebas: GraphQL (POST http://localhost:8080/graphql)
### Nota para Postman: Selecciona el método POST, usa la URL del endpoint y en la pestaña Body elige la opción nativa GraphQL. Pega estos bloques en la caja de texto izquierda (Query).

## Operación: Create (Mutation)
### GraphQL
mutation {
  create(
    nombre: "Reloj de Titanio Saavedra"
    descripcion: "Modelo Elegance de 38mm con caja sostenible"
    precio: 450.00
  ) {
    id
    nombre
    precio
  }
}
## Operación: Get (Existente - Query)
### GraphQL
query {
  get(id: "1") {
    nombre
    descripcion
    precio
  }
}
## Operación: Get (Manejo de Errores - Inexistente - Query)
### GraphQL
query {
  get(id: "999") {
    id
    nombre
  }
}
(Debe retornar data: null y un arreglo de errors con tu mensaje de validación).
## Operación: Update (Mutation)
### GraphQL
mutation {
  update(
    id: "1"
    nombre: "Reloj de Titanio Saavedra Elegance v2"
    descripcion: "Caja de titanio de 38mm y calibre automático modificado"
    precio: 480.00
  ) {
    id
    nombre
    precio
  }
}
## Operación: List (Query)
### GraphQL
query {
  list {
    id
    nombre
    precio
  }
}
## Operación: Delete (Mutation)
### GraphQL
mutation {
  delete(id: "1") {
    id
    nombre
  }
}