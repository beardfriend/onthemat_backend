# Ontheamt API Server

<div style="display:flex;">
   <img src="https://img.shields.io/badge/Go-gray?style=flat&logo=Go&logoColor=00ADD8"/>
	<img src="https://img.shields.io/badge/fiber-gray?style=flat"/>
	<img src="https://img.shields.io/badge/entGo-gray?style=flat"/>
  <img src="https://img.shields.io/badge/postgreSQL-gray?style=flat&logo=PostgreSQL&logoColor=4169E1"/>
  <img src="https://img.shields.io/badge/redis-gray?style=flat&logo=Redis&logoColor=DC382D"/>
    <img src="https://img.shields.io/badge/elastic-gray?style=flat&logo=ElasticSearch&logoColor=005571"/>
</div>
<br/>
요가 채용서비스를 위한 API 서버 


## 서버 아키텍처

![onthemat](https://user-images.githubusercontent.com/97140962/216040777-6e06651a-2ae8-4b61-ada7-4c7b9f14ab7d.png)

## 특징

-   fiber & postgreSQL을 활용한 높은 퍼포먼스의 서버
-   엘라스틱을 활용한 검색어 자동완성 시스템
-   RFC 3986를 준수한 REST API
-   테스트코드 작성 (unit, mock service logic)
-   jwt토큰인증, oAuth
-   코드를 생성하는 ORM을 사용하여 컴파일 단계에서 에러 검출.
-   apidocs를 활용하여 API 문서 작성
-   에러원인을 상세하게 안내하는 에러전달 시스템

# 시작하기

🙏

## 프로덕션

API 명세서 : http://13.125.48.238:3000/

## 로컬

**환경**

go 1.19.3  

docker  

redis    

postgresSQL

```bash

#----------- 준비중 -----------

# make database
make docker_postgres_dev

# orm으로 작성된 repository코드 생성
make generate

# 테스트용 데이터베이스 서버 올리기
make test_up

# 테스트용 데이터베이스 서버 내리기
make test_down

# seeding
make seed

# Mock코드 생성
make mockery

# 실행하기
make run

#document
make apidoc

```

### 테스트영상(데모)

https://user-images.githubusercontent.com/97140962/210781527-ece2bf08-1c96-47fa-b04b-9a58e513239e.mov

# 기타

-   fiber, entGo, elasticSearch를 처음 사용하는 프로젝트입니다.
-   학습목적으로 작성된 코드입니다.
-   모듈간의 의존성을 낮추기 위해 노력하였습니다.
-   참고: 실행파일은 cmd에 있습니다. DTO는 transport에 있습니다.
