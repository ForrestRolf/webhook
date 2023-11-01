<script setup>
import {LeftCircleOutlined} from "@ant-design/icons-vue"
import {useRouter} from "vue-router";
import TriggerGroup from "../components/triggers/TriggerGroup.vue";
import {ref, watch} from "vue";
import Action from "../components/actions/Action.vue";

const router = useRouter()

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

const triggers = ref({
    "and": [
        {
            "match": {
                "parameter": {
                    "source": "payload",
                    "name": "X-Hub-Signature"
                },
                "type": "regex",
                "regex": "mysecret"
            }
        },
        {
            "or": [
                {
                    "match":
                        {
                            "parameter":
                                {
                                    "source": "payload",
                                    "name": "ref"
                                },
                            "type": "value",
                            "value": "refs/heads/master"
                        }
                },
                {
                    "match":
                        {
                            "parameter":
                                {
                                    "source": "payload",
                                    "name": "X-GitHub-Event"
                                },
                            "type": "value",
                            "value": "ping"
                        }
                }
            ]
        },
        {
            "match": {
                "parameter": {
                    "source": "payload",
                    "name": "X-Hub-Signature"
                },
                "type": "value",
                "value": "mysecret"
            }
        },
    ]
})

const actions = ref([
    {
        "driver": "shell",
    },
    {
        "driver": "http",
    }
])

watch(triggers, () => {
    console.log(triggers.value)
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
                        <a-input :disabled="!stepVisible('basic')" ></a-input>
                    </a-form-item>
                    <a-form-item label="Description">
                        <a-textarea :rows="4" :disabled="!stepVisible('basic')"></a-textarea>
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
                <a-button type="primary" class="step-btn" v-show="stepVisible('actions')">Save</a-button>
            </a-space>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
