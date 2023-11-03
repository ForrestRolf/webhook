<script setup>
import {LeftCircleOutlined, PlusCircleOutlined} from "@ant-design/icons-vue"
import {useRouter} from "vue-router";
import TriggerGroup from "../components/triggers/TriggerGroup.vue";
import {reactive, ref, watch} from "vue";
import Action from "../components/actions/Action.vue";
import useAxios from "../support/axios.js";
import PassArgument from "../components/PassArgument.vue";
import useMessage from "../support/message.js";

const router = useRouter()
const {httpPost} = useAxios()
const {successMessage} = useMessage()

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
    pass_arguments_to_action: []
})

const triggers = ref({
    "and": [
        {
            "match": {
                "parameter": {
                    "source": "payload",
                    "name": ""
                },
                "type": "value",
                "value": ""
            }
        }
    ]
})

const actions = ref([
    {
        "driver": "shell",
        "attributes": {}
    }
])

const argument = ref([

])
const loading = ref(false)

const handleSave = () => {
    loading.value = true
    httpPost("/webhook").withBody(formState).exec().then(({payload, meta}) => {
        loading.value = false
        successMessage("Webhook created", payload.name).show()
        backToHome()
    }).catch(err => {
        loading.value = false
    })
}

const addArgument = () => {
    argument.value.push({
        source: "payload",
        value: "",
        envname: "",
    })
}
const handleRemoveArgument = (i) => {
    argument.value.splice(i, 1)
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
    formState.pass_arguments_to_action = argument.value
}, {
    deep: true
})
</script>

<template>
    <a-row :gutter="[12, 12]">
        <a-col :span="24">
            <a-space align="center">
                <a-button type="text" @click="backToHome">
                    <LeftCircleOutlined class="mini-back-btn"/>
                </a-button>
                <a-typography-title :level="3" class="no-margin">New Webhook</a-typography-title>
            </a-space>
        </a-col>
        <a-col :span="6">
            <a-card title="Basic">
                <a-form :label-col="{ span: 6 }">
                    <a-form-item label="Name">
                        <a-input :disabled="!stepVisible('basic')" v-model:value="formState.name" ></a-input>
                    </a-form-item>
                    <a-form-item label="Description">
                        <a-textarea :rows="4" :disabled="!stepVisible('basic')" v-model:value="formState.description"></a-textarea>
                    </a-form-item>
                </a-form>
            </a-card>
            <a-button type="primary" class="float-rgt step-btn" v-show="stepVisible('basic')" @click="handleNextStep">
                Next
            </a-button>
        </a-col>
        <a-col :span="10">
            <a-card title="Trigger rules">
                <TriggerGroup v-model:triggers="triggers" :disabled="!stepVisible('triggers')">

                </TriggerGroup>
            </a-card>
            <br/>
            <a-card title="Arguments">
                <a-row :gutter="[12, 12]">
                    <a-col :span="24" v-for="(arg, i) in argument">
                        <PassArgument v-model:argument="argument[i]" @remove="handleRemoveArgument(i)"></PassArgument>
                    </a-col>
                    <a-col :span="24">
                        <a-button @click="addArgument">
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
                <Action v-model:actions="actions" :disabled="!stepVisible('actions')"></Action>
            </a-card>
            <a-space class="float-rgt">
                <a-button type="default" class="step-btn" v-show="stepVisible('actions')" @click="handlePrevStep">Prev
                </a-button>
                <a-button type="primary" class="step-btn" v-show="stepVisible('actions')" @click="handleSave" :loading="loading">Save</a-button>
            </a-space>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
