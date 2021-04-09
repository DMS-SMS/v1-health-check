# **SMS-Health-Check**

> ## 대덕소프트웨어마이스터고등학교 [**SMS(학교 지원 시스템)**](https://github.com/DMS-SMS/v1-api-gateway) 서버 **Health Check 서비스**

<br>

---
## **INDEX**
### [**1. SMS Health Check란?**](#SMS-Health-Check란?)
### [**2. Health Check 종류**](#Health-Check-종류)
### [**3. Clean한 코드에 대해**](#Clean한-코드에-대해)
### [**4. 패키지 종류 및 구조**](#패키지-종류-및-구조)
### [**5. 패키지 의존성 그래프**](#패키지-의존성-그래프)

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
> ### **Clean한 코드의 기준은 무엇일까요? 제가 생각하는 기준은 다음과 같습니다.**
### 1. **의존적인 관계**에서 서로의 **계층**을 명확하게 **분리**하였고, 그 관계가 **느슨하게 결합**되었는가?
#### `-> 모든 의존 관계 추상화 완료시, mocking을 이용한 business logic 단위 테스트 가능`
- 기능에 따라 계층을 분리하고, **분리된 계층별로 패키지**를 생성하여 해당 패키지에 구현
- 직접적인 **타입의 명시**가 아닌 **인터페이스 추상화**를 통해 느슨한 상하위 계층 간의 의존 관계 형성
- 상위 계층에서, 하위 계층의 **처리 방식을 모른채로** 데이터 처리의 **책임을 빌려주는** 구조
- 따라서 계층끼리의 의존성을 편리하게 관리하기 위해서 **모든 의존성 주입**은 **main에서 발생**
### 2. 하위 계층과의 **의존 관계를 표현하는 추상화**(인터페이스)를 **상위 계층이 소유**하고 있는가?
#### `-> DIP 패턴 적용 완료시, domain을 제외한 모든 패키지의 임포트가 main 패키지에서만 발생`
- 만약 인터페이스를 **하위 계층이 소유**하고 있다면, 여전히 **하위 계층에 명시적으로 결합**된 상태
- 따라서 인터페이스 소유권을 **사용하는 계층으로 옮김**으로써, 하위 계층과의 **명시적인 결합**을 완전히 **끊을** 수 있다.
- SOLID 중 **DIP(의존성 역전)** 원칙으로, **사용할 메서드들만 추상화**하여 의존성을 생성할 수 있다는 장점 또한 존재한다.
- 예외) **domain model 관련 패키지**(repo, use case)들에 대한 추상화는 **[domain](https://github.com/DMS-SMS/v1-health-check/tree/develop/domain) 패키지에 묶어서 관리**

<br>

---
## **패키지 종류 및 구조**
> ### **프로젝트를 구성하는 패키지들을 크게 세 가지 종류로 나누어 설명합니다.**
### 1. **App**
- [**app**](https://github.com/DMS-SMS/v1-health-check/tree/develop/app)
    - **main function**을 가지고 있는 **main package**로, Health Check를 실행시키는 시작점
    - 모든 **의존성 객체 생성 및 주입**이 여기서 일어나며, [**domain 패키지**](https://github.com/DMS-SMS/v1-health-check/tree/develop/domain)를 제외한 다른 패키지를 명시적으로 import하는 유일한 패키지
- [**app/config**](https://github.com/DMS-SMS/v1-health-check/tree/develop/app/config)
    - app(main) 패키지에서 사용하는 **config value**들을 **관리**하고 **반환**하는 패키지
    - **싱글톤 패턴**으로 구현되어 있으며, **environment variable** 또는 **fixed value** 반환
    - 특정 **인터페이스의 구현체**가 아니라, 단순히 app 패키지에서 **명시적**으로 불러와서 사용하는 객체이다.

### 2. **Domain**
- [**domain**](https://github.com/DMS-SMS/v1-health-check/tree/develop/domain)
    - 특정 **domain**에서 사용할 **model 정의**와 그와 **연관된 계층**들(repo, use case)을 **추상화**하는 것이 필요
    - 따라서 이와 관련된 것을 **해당 패키지**에서 **묶어서 관리**하며, 추상화에 대한 **구현**은 **domain 이름으로 된 패키지** 내부에서 진행한다.
    - 결론적으로 도메인에 대한 **model struct, repository와 usecase interface**를 정의한다.
    - 현재 추상화 및 구현이 완료된 domain에는 **syscheck**와 **srvcheck**가 있다.
- [**syscheck**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck)
    - **system check** 기능의 **domain에 대한 추상화**를 실제로 **구현** 하는 패키지 
    - **domain 패키지**에서 정의한 **추상화**에 의존하고 있으며, **네 가지의 하위 패키지**로 구성되어있다.
    - [**repository**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck/repository)
        - domain 패키지에서 **추상화**된 **system check** 관련 **repository**들을 구현하는 패키지
        - domain 패키지에 정의된 **model struct에 의존**하고 있으며, 데이터를 **명령 혹은 조회**하는 기능의 계층이다.
        - 현재로써는 **elasticsearch**를 저장소로 사용하는 구현체만 존재한다.
    - [**usecase**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck/usecase)
        - domain 패키지에서 **추상화**된 **system check** 관련 **usecase**들을 구현하는 패키지
        - domain 패키지에 정의된 **repository 추상화에 의존**하고 있으며, 실질적인 **business logic**을 처리하는 기능의 계층이다.
        - 또한 **외부 서비스들**(docker, slack, etc..)에 대해서도 **추상화에 의존**하고 있으며, **해당 추상화에 대한 소유권**은 해당 패키지가 가지고 있다.
        - 이러한 추상화에 대한 구현체는 **Agent 관련 패키지**에서 확인할 수 있다.
    - [**delivery**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck/delivery)
        - 특정 API로부터 들어온 데이터를 **usecase layer으로 전달**하는 기능의 계층
        - 따라서, domain 패키지에 정의된 **usecase 추상화에 의존**하고 있다.
        - 해당 프로젝트 내에서 **최상위 계층**으로, 어떠한 추상화에 대한 구현체가 아니다.
    - [**config**](https://github.com/DMS-SMS/v1-health-check/tree/develop/syscheck/config)
        - **syscheck의 모든 하위 패키지**에서 사용하는 **config value**들을 **관리**하고 **반환**하는 패키지
        - **싱글톤 패턴**으로 구현되어 있으며, [**config.yaml**](https://github.com/DMS-SMS/v1-health-check/blob/develop/config.yaml) 파일에 설정된 값 또는 기본 값 반환
        - repository, usecase, delivery 패키지에서 **추상화된 인터페이스**들을 **모두 구현**하고 있음.
    - 같은 도메인 내에서도 기능들끼리의 연관성을 없애기 위해, **모든 기능들에 대한 추상화와 구현체들이 서로 다른 타입으로 분리되어있다.**
- [**srvcheck**](https://github.com/DMS-SMS/v1-health-check/tree/develop/srvcheck)
    - syscheck 패키지와 비슷하게, **service check** 기능의 domain에 대한 **추상화**를 **구현**하는 패키지이다.
    - syscheck 패키지와 하위 구성 또한 동일하지만, 서로 간의 **결합**이 전혀 **존재하지 않다.**
### 3. **Agent**
- [**consul**](https://github.com/DMS-SMS/v1-health-check/tree/develop/consul)
- 

<!-- 
![godepgraph1](https://user-images.githubusercontent.com/48676834/113800510-08960180-9792-11eb-8c5d-a5650ab0799b.png)

![godepgraph2](https://user-images.githubusercontent.com/48676834/113800517-0df34c00-9792-11eb-90ab-d048f3b847a1.png) -->
