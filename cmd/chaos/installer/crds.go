package installer

// Taken from https://litmuschaos.github.io/litmus/litmus-operator-v1.7.0.yaml
//
// TODO(giri): Need to convert this to CRD code

var chaosEngineCRD = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: chaosengines.litmuschaos.io
spec:
  group: litmuschaos.io
  names:
    kind: ChaosEngine
    listKind: ChaosEngineList
    plural: chaosengines
    singular: chaosengine
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
          properties:
            monitoring:
              type: boolean
            jobCleanUpPolicy:
              type: string
              pattern: ^(delete|retain)$
              # alternate ways to do this in case of complex pattern matches
              #oneOf:
              #  - pattern: '^delete$'
              #  - pattern: '^retain$'
            annotationCheck:
              type: string
              pattern: ^(true|false)$
            appinfo:
              type: object
              properties:
                appkind:
                  type: string
                  pattern: ^(deployment|statefulset|daemonset|deploymentconfig)$
                applabel:
                  pattern: ([a-z0-9A-Z_\.-/]+)=([a-z0-9A-Z_\.-/_]+)
                  type: string
                appns:
                  type: string
            auxiliaryAppInfo:
              type: string
            engineState:
              type: string
              pattern: ^(active|stop|initialized|stopped)$
            chaosServiceAccount:
              type: string
            components:
              type: object
              properties:
                runner:
                  type: object
                  properties:
                    image:
                      type: string
                    type:
                      type: string
                      pattern: ^(go)$
                    runnerannotation:
                      type: object
                      additionalProperties:
                        type: string
                        properties:
                          key:
                            type: string
                            minLength: 1
                            allowEmptyValue: false
                          value:
                            type: string
                            minLength: 1
                            allowEmptyValue: false
            experiments:
              type: array
              items:
                type: object
                properties:
                  name:
                    type: string
                  spec:
                    type: object
                    properties:
                      k8sProbe:
                        type: array
                        items:
                          type: object
                          properties:
                            name:
                              type: string
                            inputs:
                              type: object
                              properties:
                                command:
                                  type: object
                                  properties:
                                    group:
                                      type: string
                                    version:
                                      type: string
                                    resource:
                                      type: string
                                    namespace:
                                      type: string
                                    fieldSelector:
                                      type: string
                                expectedResult:
                                  type: string
                            runProperties:
                              type: object
                              properties:
                                probeTimeout:
                                  type: integer
                                interval:
                                  type: integer
                                retry: 
                                  type: integer
                            mode:
                              type: string
                      cmdProbe:
                        type: array
                        items:
                          type: object
                          properties:
                            name:
                              type: string
                            inputs:
                              type: object
                              properties:
                                command:
                                  type: string
                                expectedResult:
                                  type: string
                                source:
                                  type: string
                            runProperties:
                              type: object
                              properties:
                                probeTimeout:
                                  type: integer
                                interval:
                                  type: integer
                                retry: 
                                  type: integer
                            mode:
                              type: string
                      httpProbe:
                        type: array
                        items:
                          type: object
                          properties:
                            name:
                              type: string
                            inputs:
                              type: object
                              properties:
                                url:
                                  type: string
                                expectedResponseCode:
                                  type: string
                            runProperties:
                              type: object
                              properties:
                                probeTimeout:
                                  type: integer
                                interval:
                                  type: integer
                                retry: 
                                  type: integer
                            mode:
                              type: string
                      components:
                        type: object
                        properties:
                          statusCheckTimeouts:
                            type: object
                            properties:
                              delay:
                                type: integer
                              timeout:
                                type: integer
                          nodeSelector:
                            type: object
                            minLength: 1
                          experimentImage:
                            type: string
                          env:
                            type: array
                            items:
                              type: object
                              properties:
                                name:
                                  type: string
                                value:
                                  type: string
                          configMaps:
                            type: array
                            items:
                              type: object
                              properties:
                                name:
                                  type: string
                                mountPath:
                                  type: string
                          secrets:
                            type: array
                            items:
                              type: object
                              properties:
                                name:
                                  type: string
                                mountPath:
                                  type: string
                          experimentannotation:
                            type: object
                            additionalProperties:
                              type: string
                              properties:
                                key:
                                  type: string
                                  minLength: 1
                                  allowEmptyValue: false
                                value:
                                  type: string
                                  minLength: 1
                                  allowEmptyValue: false


        status:
          type: object
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
`

var chaosExperimentCRD = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: chaosexperiments.litmuschaos.io
spec:
  group: litmuschaos.io
  names:
    kind: ChaosExperiment
    listKind: ChaosExperimentList
    plural: chaosexperiments
    singular: chaosexperiment
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
          properties:
            definition:
              type: object
              properties:
                args:
                  type: array
                  items:
                    type: string
                command:
                  type: array
                  items:
                    type: string
                env:
                  type: array
                  items:
                    type: object
                    properties:
                      name:
                        type: string
                      value:
                        type: string
                image:
                  type: string
                labels:
                  type: object
                  properties:
                    name:
                      type: string
                scope:
                  type: string
                  pattern: ^(Namespaced|Cluster)$
                permissions:
                  type: array
                  items:
                    type: object
                    minProperties: 3
                    required:
                      - apiGroups
                      - resources
                      - verbs
                    properties:
                      apiGroups:
                        type: array
                        items:
                          type: string
                      resources:
                        type: array
                        items:
                          type: string
                      verbs:
                        type: array
                        items:
                          type: string
                      resourceNames:
                        type: array
                        items:
                          type: string
                      nonResourceURLs:
                        type: array
                        items:
                          type: string
                configmaps:
                  type: array
                  items:
                    type: object
                    minProperties: 2
                    properties:
                      name:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                      mountPath:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                secrets:
                  type: array
                  items:
                    type: object
                    minProperties: 2
                    properties:
                      name:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                      mountPath:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                hostFileVolumes:
                  type: array
                  items:
                    type: object
                    minProperties: 3
                    properties:
                      name:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                      mountPath:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                      nodePath:
                        type: string
                        allowEmptyValue: false
                        minLength: 1
                securityContext:
                  type: object
                hostPID:
                  type: boolean
        status:
          type: object
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
`

var chaosResultCRD = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: chaosresults.litmuschaos.io
spec:
  group: litmuschaos.io
  names:
    kind: ChaosResult
    listKind: ChaosResultList
    plural: chaosresults
    singular: chaosresult
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
        status:
          type: object
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
`
