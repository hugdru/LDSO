### Backend ###
  * *Language* go
   * [pressly/chi](https://github.com/pressly/chi) - Ajuda a configurar as rotas de maneira organizada
   * [rs/cors](https://github.com/rs/cors) - Biblioteca [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS) que permite à frontend fazer pedidos à API
   * [elithrar/simple-scrypt](https://github.com/elithrar/simple-scrypt) - Hashing e Salting seguro de passwords
   * [jmoiron/sqlx](https://github.com/jmoiron/sqlx) - Facilita o uso de sql, extende a biblioteca principal
   * [guregu/null.v3](https://github.com/guregu/null) - Melhora o tratamento de nulls da base de dados por parte do go
   * [alexedwards/scs](https://github.com/alexedwards/scs) - Gere a sessão de um utilizador e permite guardar em diversos sítios
   * [alexedwards/scs/engine/redisstore](https://github.com/alexedwards/scs/tree/master/engine/redisstore) - Biblioteca necessária para guardar a sessão no redis
  * *Database Bulk* [postgresql](https://www.postgresql.org/) - Guardamos todos os dados da aplicação incluindo imagens (bytea) aqui
  * *Session Management* [redis](https://redis.io/) - Usado para guardar os dados da sessão de cada utilizar
  * *Build Tool* [gb](https://getgb.io/) - Ferramenta que permite guardar o código go em modo projeto vs workspace(default do go, tipo eclipse). Também facilita a gestão de bibliotecas: instalação, atualização, controlo de versão, etc. Ver como o go organiza o código por default em [go organization](https://golang.org/doc/code.html#Organization). E posteriormente comparar com a estrutura de projeto imposta pelo gb.

### Frontend ###
  * *Framework* [angular](https://angular.io/) 2 with typescript and angular-cli - Usado para desenhar a single page app (SPA). Faz pedidos à API escrita em go. Usámos o angular-cli porque facilita a criação da estrutura base, a atualização, execução, etc.
  * *Responsive Framework* [bootstrap](http://getbootstrap.com/) - Para facilitar e uniformizar o aspecto do site

### Tools ###
  * *IDE* webstorm<br>
  * *Continuous Deployment* [docker](https://www.docker.com/) - usado para criar imagens que posteriormente devem ser postas no serviço cloud. Usámos o digitalocean durante o desenvolvimento, existem outros, por exemplo: amazon ec2.
  * *Continuous Deployment Aids* [docker-compose](https://docs.docker.com/compose/) - pega nas imagens e cria uma rede em que os contendores (imagens em execução) podem comunicar entre si
  * *Continuous Integration Tool* [travis-ci](https://travis-ci.org/) - Usado para fazer alguns testes automaticamente

### Physical Architecture ###
![Physical Architecture Image](/docs/physical.png)

### File Structure ###
 * apiTests/
  * specs/ - ficheiros com os resultados esperados
  * apiTests.sh - testes curl que comparam as especificações com o resultado do pedido
 * backend/
  * vendor/
    * manifest - ficheiro do gb que guarda a informação das bibliotecas usadas
  * src/server/
    * datastore/ - modelos, validação do modelo, sql e conexões ao postgresql e ao redis. Funciona como um wrapper de queries, para não estarmos a fazer queries iguais ou validações repetidas em cada handler.
     * datastore.go - ficheiro com as conexões
       * generators/ - geradores de partes de código sql para ins, filtros, etc com bindings (prepared statements)
    * handler/ - rotas e ficheiros auxiliares. Com o chi é possível dividir em partes a API (.Route), isto melhora a legibilidade e permite agrupar código comum através do contexto(novo no go 1.7) que é passado nos middlewares(.Use, ).
      * handler.go - ficheiro com as subrotas
      * sessionData/ - contém funções de auxílio para lidar com a sessão
      * helpers/
        * datetime.go - maneira universal de lidar com datas
        * functions.go - misc de funções úteis
        * images.go - funções que lidam com as imagens
        * decorators/ - decoradores que envolvem os handlers e que fazem validação, etc antes de chegar ao código do handler
 * frontend/ ficheiros do angular2
   * frontend/src/environments/ - ficheiros com os urls da API dependendo se staging ou prod. Usado no Docker, chama ng build --env="$BUILD" . O ficheiro environment.ts é substituído pelo angular2 dependendo de "$BUILD".
 * database/
   * sql/
     * \*.sql - ficheiros sql ordenados, ordem de execução. Estes ficheiros lidam com o schema. Creates, Indexes, Triggers, etc.
     * examples/ - ficheiros sql ordenados, ordem de execução. Inserts e updates exemplares.
   * uml/classdiagram.txt - ficheiro plantuml com o class diagram da base de dados. Ligeiramente desatualizado, comparar com 1.0create.sql.
   * allin.sh - script que concatena tudo o que estiver na pasta sql/ num só ficheiro por ordem e de forma recursiva.
   * pgr - script que faz wrapping ao postgresql, e que permite iniciar um cluster local rapidamente.
 * docker/
   * common/ - ficheiros docker comuns a todas as partes
   * coreos/ - ficheiros de configuração da distribuição linux coreos
   * development/ - ficheiros docker e docker-compose associados à máquina de cada programador
   * production/ - ficheiros docker e docker-compose associados ao servidor de production
   * staging/ - ficheiros docker e docker-compose associados ao servidor de staging
   * templates/ - ficheiros usados na criação de imagens do docker que podem ser reutilizados
   * exclude_context - ficheiro usado pelo script dkr que previne certos ficheiros serem enviados no contexto do docker quando uma imagem é criada
 * setup/
   * debian_based_tools - ficheiro de ajuda de instalação de ferramentas no debian, encontra-se desatualizado
 * .travis.yml - ficheiro do travis que efetua testes automaticamente
 * dkr - script que auxilia a execução, substituição, criação, etc das imagens do docker e que se pretendido as coloca para os servidores

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

### How to run the project ###
   1. Run an updated linux distribution, for instance [ubuntu](https://www.ubuntu.com/).
   2. Install [go >= 1.7](https://golang.org/doc/install)
   3. Install gb after configuring go paths and PATH, go get github.com/constabulary/gb/...
   4. Install [nodejs and npm](https://nodejs.org/en/download/)
   5. Add the npm bin folder to PATH
   6. Install [angular-cli](https://github.com/angular/angular-cli) globally, npm install -g angular-cli
   7. Install [redis](https://redis.io/download) and [postgresql](https://www.postgresql.org/download/)
   8. Go to the backend folder and do, gb vendor restore , gets the backend dependencies
   9. Go to the frontend folder and do, npm install , gets the frontend dependencies
   10. Add the following to /etc/hosts . postgres and redis hosts are important because these are hardcoded in the backend and used in docker. api.lp4adev.tk can be set in the frontend environments as explained previously.
       1. 127.0.0.1 go1
       2. 127.0.0.1 gon
       3. 127.0.0.1 postgres
       4. 127.0.0.1 redis
       5. 127.0.0.1 lp4adev.tk
       6. 127.0.0.1 www.lp4adev.tk
       7. 127.0.0.1 api.lp4adev.tk
   11. Check if redis is already running in your computer if not run it. systemctl status|start|stop|enable|disable program-name . ps aux | grep redis
   12. Go to the frontend folder and do ng server
   13. Go to the database folder and do ./pgr initddlexamples
   14. Go to the backend folder and do gb build server && ./bin/server
   15. Insert lp4adev.tk:4200 in a browser.

### How to build the images and send to the server ###
   1. Run an updated linux distribution, for instance [ubuntu](https://www.ubuntu.com/).
   2. Install [docker](https://docs.docker.com/engine/installation/)
   3. (Optional) Install [docker-compose](https://docs.docker.com/compose/install/)
   4. Make sure the docker daemon is running.
   5. You might need to add your user to the docker group
   6. Run ./dkr redo prod . To turn the remote containers off, build all the prod images, clean all remote images, push all  the prod images, push the prod docker-compose file, turn the remote containers on.
