### Backend ###
  * *Language* go
   * [pressly/chi](https://github.com/pressly/chi) - Ajuda a configurar as rotas de maneira organizada
   * [rs/cors](https://github.com/rs/cors) - Biblioteca [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS) que permite aceitar pedidos da frontend
   * [elithrar/simple-scrypt](https://github.com/elithrar/simple-scrypt) - Hashing e Salting seguro de passwords
   * [jmoiron/sqlx](https://github.com/jmoiron/sqlx) - Facilita o uso de sql, extende a biblioteca principal
   * [guregu/null.v3](https://github.com/guregu/null) - Melhora o tratamento de nulls da base de dados por parte do go
   * [alexedwards/scs](https://github.com/alexedwards/scs) - Gere a sessão de um utilizador e permite guardar em diversos sítios
   * [alexedwards/scs/engine/redisstore](https://github.com/alexedwards/scs/tree/master/engine/redisstore) - Biblioteca necessária para guardar a sessão no redis
  * *Database Bulk* [postgresql](https://www.postgresql.org/) - Guardamos todos os dados da aplicação incluindo imagens (bytea) aqui
  * *Session Management* [redis](https://redis.io/) - Usado para guardar os dados da sessão de cada utilizar
  * *Build Tool* [gb](https://getgb.io/) - Ferramenta que permite guardar o código go em modo projeto vs workspace(default do go) tipo eclipse. Também facilita a gestão de bibliotecas: instalação, atualização, controlo de versão, etc. Ver como o go organiza o código por default em [go organization](https://golang.org/doc/code.html#Organization). E posteriormente comparar com a estrutura de projeto imposta pelo gb.

### Frontend ###
  * *Framework* [angular](https://angular.io/) 2 with typescript and angular-cli - Usado para desenhar a single page app (SPA). Faz pedidos à API escrita em go. Usámos o angular-cli porque facilita a criação da estrutura base, a atualização, execução, etc.
  * *Responsive Framework* [bootstrap](http://getbootstrap.com/) - Para facilitar e uniformizar o aspecto do site

### Tools ###
  * *IDE* webstorm<br>
  * *Continuous Deployment* [docker](https://www.docker.com/) - usado para criar imagens que posteriormente devem ser postas no serviço cloud. Usámos o digitalocean durante o desenvolvimento, existem outros, por exemplo: amazon ec2.
  * *Continuous Deployment Aids* [docker-compose](https://docs.docker.com/compose/) - pega nas imagens e cria uma rede em que os contendores (imagens em execução) podem comunicar entre si
  * *Continuous Integration Tool* [travis-ci](https://travis-ci.org/) - Usado para fazer alguns testes automaticamente

### Resources ###

#### Go ####
  * https://tour.golang.org/list
  * https://golang.org/doc/articles/wiki/
  * https://golang.org/doc/code.html
  * https://getgb.io/docs/usage/
  * https://godoc.org/github.com/constabulary/gb/cmd/gb-vendor
  * https://getgb.io/

#### Angular 2 ####
  * https://angular.io/docs/ts/latest/quickstart.html

#### HTTP ####
  * https://code.tutsplus.com/tutorials/http-the-protocol-every-web-developer-must-know-part-1--net-31177
  * https://code.tutsplus.com/tutorials/http-the-protocol-every-web-developer-must-know-part-2--net-31155
