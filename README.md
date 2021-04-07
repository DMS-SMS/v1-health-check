# **SMS-Health-Check**

> ### 대덕소프트웨어마이스터고등학교 [**SMS(학교 지원 시스템)**](https://github.com/DMS-SMS) 서버 **Health Check 서비스**

<br>

---
## **INDEX**
### [**1. SMS Health Check란?**](#SMS-Health-Check란?)
### [**2. Health Check 종류**](#Health-Check-종류)
### [**3. 패키지 의존 흐름 및 그래프**](#패키지-의존성-흐름-및-그래프)
### [**4. 프로젝트 구조**](#프로젝트-구조)
### [**5. 부가 기능**](#부가-기능)

<br>

---
## **SMS Health Check란?**
- **SMS**(School Management System)는 전공 동아리 [**DMS**](https://github.com/DSM-DMS)에서 개발하여 현재 운영하고 있는 **학교 지원 시스템**입니다.
- **SMS Health Check**는 SMS 서버를 **운영중인 환경**과 서버를 **구성중인 여러 서비스**들의 상태를 **주기적으로 관리**하는 플러그인 식의 서비스입니다.
- 특정 서비스의 **상태 확인 결과**가 특정 기준을 통해 **정상적이지 않다고 판단**이 되면 해당 서비스의 **상태 회복을 위한 작업을 수행**합니다.

<br>

---
## **Health Check 종류**
#### **`모든 Health Check는 수행 결과를 Elasticsearch에 저장하여 관리합니다.`**
### 1. [**System Check**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck)
- **disk check**
    - **디스크 사용 용량 특정 수치 초과** 시 알람 발행 후 **Docker Prune 실행**
    - 디스크 사용의 주요 원인에는 로그 및 DB 데이터 또한 존재하지만,  Docker Pruning만 실행
- **memory check**
    - **총 메모리 사용량 특정 수치 초과** 시 알람 발행 후 **메모리 과다 사용 프로세스 재부팅**
    - 프로세스 메모리 사용 조회 및 재부팅은 **Docker Engine API**를 통해 수행

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
    - 또한, 서비스별로 **등록된 노드가 존재하지 않는** 경우 알람 발행 후 **해당 서비스 재부팅**
    - 작동되지 않는 노드인지는 해당 노드와 **gRPC 연결 시도**를 통해 판별
    - 노드 부재시에는, **서비스 재부팅**을 함으로써 재시작 시점에 **스스로 노드를 등록**하게 됨

<!-- 
![godepgraph1](https://user-images.githubusercontent.com/48676834/113800510-08960180-9792-11eb-8c5d-a5650ab0799b.png)

![godepgraph2](https://user-images.githubusercontent.com/48676834/113800517-0df34c00-9792-11eb-90ab-d048f3b847a1.png) -->
