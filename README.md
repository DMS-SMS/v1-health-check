# **SMS-Health-Check**

> ### 대덕소프트웨어마이스터고등학교 [**SMS(학교 지원 시스템)**](https://github.com/DMS-SMS/v1-api-gateway) 서버 **Health Check 서비스**

<br>

---
## **INDEX**
### [**1. SMS Health Check란?**](#SMS-Health-Check란?)
### [**2. Health Check 종류**](#Health-Check-종류)
### [**3. Clean한 코드에 대해**](#Clean한-코드에-대해)
### [**4. 패키지 의존 흐름 및 그래프**](#패키지-의존-흐름-및-그래프)
### [**5. 프로젝트 구조**](#프로젝트-구조)
### [**6. 부가 기능**](#부가-기능)

<br>

---
## **SMS Health Check란?**
- [**SMS**](https://github.com/DMS-SMS/v1-api-gateway)(School Management System)는 대덕소프트웨어마이스터고등학교 전공 동아리 [**DMS**](https://github.com/DSM-DMS)에서 개발하여 현재 운영하고 있는 **학교 지원 시스템**입니다.

- [**SMS Health Check**](https://github.com/DMS-SMS/v1-health-check)는 SMS 서버를 **운영중인 환경**과 서버를 **구성중인 여러 서비스**들의 상태를 **주기적으로 관리**하는 플러그인 식의 서비스입니다.
- 특정 서비스의 **상태 확인 결과**가 특정 기준을 통해 **정상적이지 않다고 판단**이 되면 해당 서비스의 **상태 회복을 위한 작업을 수행**합니다.

<br>

---
## **Health Check 종류**
> ### **모든 Health Check는 수행 결과를 Elasticsearch에 저장하여 관리합니다.**
### 1. [**System Check**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck)
- **disk check**
    - **디스크 사용 용량 특정 수치 초과** 시 알람 발행 후 **Docker Prune 실행**
    - 디스크 사용의 주요 원인에는 **로그 및 DB 데이터** 또한 존재하지만, 해당 데이터를 건드는건 **위험**하다 판단하여 Docker Prune만 실행

- **memory check**
    - **총 메모리 사용량 특정 수치 초과** 시 알람 발행 후 **메모리 과다 사용 프로세스 재부팅**
    - 프로세스 메모리 사용 조회 및 재부팅은 **Docker Engine API**를 통해 수행
    - 참고로, 서버 구성에 있어서 **부재가 발생하면 안되는 서비스**들은 재부팅하지 않음

### 2. [**Service Check**](https://github.com/DMS-SMS/v1-health-check/tree/develop/srvcheck)
- **elasticsearch check**
    - **Elasticsearch Shard 갯수 특정 수치 초과** 시 알람 발행 후 **Jaeger Index 삭제**
    - Jaeger에서 **매일 새로운 Index를 생성**하여 데이터를 저장하기에 **주기적으로 관리 필요**
    - 예전에 생성된 Index부터 삭제를 진행하며, **최근 한 달 내**에 생성된 Index는 **삭제되지 않음**

- **swarmpit check**
    - **Swarmpit App 컨테이너 메모리 사용량 특정 기준 초과** 시 알람 발행 후 **해당 서비스 재부팅**
    - Swarmpit App의 경우, **메모리 사용량이 지속적으로 증가**하기 때문에 특정 수치에 도달할 때 마다 **재부팅**이 필요함
- **consul check**
    - Consul에 **작동되지 않는 노드가 등록**되었다면 알람 발행 후 **해당 노드 등록 해제**
    - 또한, MSA 상의 서비스별로 **등록된 노드가 존재하지 않는** 경우 알람 발행 후 **해당 서비스 재부팅**
    - 작동되지 않는 노드인지는 해당 노드와 **gRPC 연결 시도**를 통해 판별
    - 노드 부재시에는, **서비스 재부팅**을 함으로써 재시작 시점에 **스스로 노드를 등록**하게 함

<br>

---
## **Clean한 코드에 대해**
> ### **Clean한 코드의 기준**은 무엇일까요? **제가 생각하는 기준**은 다음과 같습니다.
### 1. **의존적인 관계**에서 서로의 **계층**을 명확하게 **분리**하였고, 그 관계가 **느슨하게 결합**되었는가?
- 기능에 따라 계층을 분리하고, **분리된 계층별로 패키지**를 생성하여 해당 패키지에 구현

- 직접적인 **타입의 명시가 아닌 인터페이스**를 통해 느슨한 상하위 계층 간의 의존 관계 형성
- 상위 계층에서, 하위 계층의 **처리 방식을 모른채로** 데이터 처리의 **책임을 빌려주는** 구조
- 따라서 계층끼리의 의존성을 편리하게 관리하기 위해서 **모든 의존성 주입**은 **main에서 발생**

### 2. 하위 계층과의 **의존 관계를 표현하는 추상화**(인터페이스)를 **상위 계층이 소유**하고 있는가?
- 만약 인터페이스를 **하위 계층이 소유**하고 있다면, 여전히 **하위 계층에 명시적으로 결합**된 상태

- 따라서 인터페이스 소유권을 **사용하는 계층으로 옮김**으로써, 하위 계층과의 **명시적인 결합**을 완전히 **끊을** 수 있음
- 예외) **domain model 관련 패키지**(repo, ucase)들에 대한 추상화는 **[domain](https://github.com/DMS-SMS/v1-health-check/tree/develop/domain) 패키지에 묶어서 관리**

<br>

명확한 범위를 가지는 패키지  
SOLID  
단위 테스트 가능성  

<br>

---
## **패키지 의존 흐름 및 그래프**
- 해당 프로젝트는 다음과 같이 크게 3 종류의 패키지로 구성되어 있습니다.

<!-- 
![godepgraph1](https://user-images.githubusercontent.com/48676834/113800510-08960180-9792-11eb-8c5d-a5650ab0799b.png)

![godepgraph2](https://user-images.githubusercontent.com/48676834/113800517-0df34c00-9792-11eb-90ab-d048f3b847a1.png) -->
