<script setup>
import {LeftCircleOutlined, PlusCircleOutlined} from "@ant-design/icons-vue"
import {useRoute, useRouter} from "vue-router";
import TriggerGroup from "../components/triggers/TriggerGroup.vue";
import {computed, onMounted, onUnmounted, reactive, ref, watch} from "vue";
import Action from "../components/actions/Action.vue";
import useAxios from "../support/axios.js";
import PassArgument from "../components/PassArgument.vue";
import useMessage from "../support/message.js";
import emitter from "../support/emitter.js";
import _ from "lodash"

const route = useRoute()
const router = useRouter()
const {httpPost, httpGet, httpPut} = useAxios()
const {successMessage, errorMessage} = useMessage()

const currentStep = ref("basic")
const steps = ref(["basic", "triggers", "actions"])

const backToHome = () => {
    router.push({name: "home"})
}

const stepVisible = (name) => {
    return name === currentStep.value
}
const handleNextStep = () => {
    let idx = steps.value.indexOf(currentStep.value)
    currentStep.value = steps.value[idx + 1]
}
const handlePrevStep = () => {
    let idx = steps.value.indexOf(currentStep.value)
    currentStep.value = steps.value[idx - 1]
}

const formState = reactive({
    name: null,
    description: null,
    triggers: {},
    actions: [],
    passArgumentsToAction: [],
    saveRequest: ["body"]
})

const match = {
    "parameter": {
        "source": "payload",
        "name": ""
    },
    "type": "value",
    "value": "",
    "not": false
}
const triggers = ref({
    "and": [
        {
            "match": {...match}
        }
    ]
})
const hasTriggers = computed(() => {
    return Object.keys(triggers.value).length > 0
})

const actions = ref([
    {
        "driver": "shell",
        "attributes": {}
    }
])
const saveRequestOptions = ref([
    {label: 'Body', value: 'body'},
    {label: 'Header', value: 'header'},
    {label: 'Query', value: 'query'},
])

const argument = ref([])
const loading = ref(false)
const webhookId = computed(() => {
    return route.query.id
})
const isEditMode = computed(() => {
    return !!route.query.id
})

const handleSave = () => {
    loading.value = true
    httpPost("/webhook").withBody(formState).exec().then(({payload, meta}) => {
        loading.value = false
        successMessage("Webhook created", payload.name).show()
        backToHome()
    }).catch(e => {
        errorMessage("Error", e).show()
    }).finally(err => {
        loading.value = false
    })
}

const handleUpdate = () => {
    loading.value = true
    httpPut(`/webhook/${webhookId.value}`).withBody(formState).exec().then(({payload, meta}) => {
        loading.value = false
        successMessage("Webhook updated", payload.name).show()
        backToHome()
    }).catch(e => {
        errorMessage("Error", e).show()
    }).finally(err => {
        loading.value = false
    })
}

const loadWebhookDetail = () => {
    if (!isEditMode) return
    if (!webhookId.value) return
    httpGet(`/webhook/${webhookId.value}`).exec().then(({payload}) => {
        for (let f of ["name", "description", "actions", "triggers", "passArgumentsToAction", "authToken", "saveRequest"]) {
            formState[f] = payload[f]
        }
        triggers.value = payload.triggers
        actions.value = payload.actions
        argument.value = payload.passArgumentsToAction || []
    }).finally(() => {

    })
}

const addArgument = () => {
    argument.value.push({
        source: "payload",
        name: "",
        envname: "",
    })
}
const handleRemoveArgument = (i) => {
    argument.value.splice(i, 1)
}

const startWithAndRule = () => {
    if (hasTriggers.value) return
    triggers.value = {
        "and": [
            {
                "match": {...match}
            }
        ]
    }
}
const startWithOrRule = () => {
    if (hasTriggers.value) return
    triggers.value = {
        "or": [
            {
                "match": {...match}
            }
        ]
    }
}

const handleRemoveTriggerGroup = (path) => {
    if (path.length === 1) {
        triggers.value = {}
    }
    path = path.map(n => n + "")
    path.pop()

    let _triggers = _.cloneDeep(triggers.value)
    _.unset(_triggers, path)

    //Filter the null left by unset
    path.pop()
    _.set(_triggers, path, _.compact(_.get(_triggers, path)))

    triggers.value = _triggers
}

watch(triggers, () => {
    formState.triggers = triggers.value
}, {
    deep: true
})
watch(actions, () => {
    formState.actions = actions.value
}, {
    deep: true
})
watch(argument, () => {
    formState.passArgumentsToAction = argument.value
}, {
    deep: true
})

onMounted(() => {
    loadWebhookDetail()
    emitter.on("trigger-group-removed", handleRemoveTriggerGroup)
})
onUnmounted(() => {
    emitter.off("trigger-group-removed", handleRemoveTriggerGroup)
})
</script>

<template>
    <a-row :gutter="[12, 12]">
        <a-col :span="24">
            <a-space align="center">
                <a-button type="text" @click="backToHome">
                    <LeftCircleOutlined class="mini-back-btn"/>
                </a-button>
                <a-typography-title :level="3" class="no-margin">{{ isEditMode ? "Edit" : "New" }} Webhook
                </a-typography-title>
            </a-space>
        </a-col>
        <a-col :span="6">
            <a-card title="Basic">
                <a-form :label-col="{ span: 6 }" :disabled="!stepVisible('basic')">
                    <a-form-item label="Name">
                        <a-input v-model:value="formState.name"></a-input>
                    </a-form-item>
                    <a-form-item label="Description">
                        <a-textarea :rows="4"
                                    v-model:value="formState.description"></a-textarea>
                    </a-form-item>
                    <a-form-item label="Auth token"
                                 help="If it is empty, it means no authorization. Otherwise, set Authorization=hook [token] in the header.">
                        <a-input v-model:value="formState.authToken"></a-input>
                    </a-form-item>
                    <a-form-item label="Save request">
                        <a-checkbox-group v-model:value="formState.saveRequest" :options="saveRequestOptions"/>
                    </a-form-item>
                </a-form>
            </a-card>
            <a-button type="primary" class="float-rgt step-btn" v-show="stepVisible('basic')" @click="handleNextStep">
                Next
            </a-button>
        </a-col>
        <a-col :span="10">
            <a-card title="Trigger rules">
                <TriggerGroup v-model:triggers="triggers" :disabled="!stepVisible('triggers')" :path="['and']"
                              v-show="hasTriggers">

                </TriggerGroup>
                <a-space v-show="!hasTriggers">
                    <a-button @click="startWithAndRule">
                        <template #icon>
                            <PlusCircleOutlined/>
                        </template>
                        Start with AND
                    </a-button>
                    <span>OR</span>
                    <a-button @click="startWithOrRule">
                        <template #icon>
                            <PlusCircleOutlined/>
                        </template>
                        Start with OR
                    </a-button>
                </a-space>
            </a-card>
            <br/>
            <a-card title="Arguments">
                <a-row :gutter="[12, 12]">
                    <a-col :span="24" v-for="(arg, i) in argument">
                        <PassArgument v-model:argument="argument[i]" @remove="handleRemoveArgument(i)"
                                      :disabled="!stepVisible('triggers')"></PassArgument>
                    </a-col>
                    <a-col :span="24">
                        <a-button @click="addArgument" :disabled="!stepVisible('triggers')">
                            <template #icon>
                                <PlusCircleOutlined/>
                            </template>
                            Add argument
                        </a-button>
                    </a-col>
                </a-row>
            </a-card>
            <a-space class="float-rgt">
                <a-button type="default" class="step-btn" v-show="stepVisible('triggers')" @click="handlePrevStep">
                    Prev
                </a-button>
                <a-button type="primary" class="step-btn" v-show="stepVisible('triggers')" @click="handleNextStep">
                    Next
                </a-button>
            </a-space>
        </a-col>
        <a-col :span="8">
            <a-card title="Actions">
                <Action v-model:actions="actions" :arguments="argument" :disabled="!stepVisible('actions')"></Action>
            </a-card>
            <a-space class="float-rgt">
                <a-button type="default" class="step-btn" v-show="stepVisible('actions')" @click="handlePrevStep">Prev
                </a-button>
                <a-button type="primary" class="step-btn" v-show="stepVisible('actions') && !isEditMode"
                          @click="handleSave" :loading="loading">Save
                </a-button>
                <a-button type="primary" class="step-btn" v-show="stepVisible('actions') && isEditMode"
                          @click="handleUpdate" :loading="loading">Update
                </a-button>
            </a-space>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
