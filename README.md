# Chuck Norris Jokes API

## Prueba técnica

API desarrollada en Go que retorna un arreglo de 25 objetos obtenidos de la API: https://api.chucknorris.io/
mostrando unicamente ciertas propiedades y sin repetir ningun elemento del arreglo.

Se cambió la función getJokesList agregando: gorountines, WaitGroup y Mutex para el manejo de la concurrencia, logrando asi reducir el tiempo de respuesta de la API
