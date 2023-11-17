<script setup>
import {computed, reactive, ref} from "vue";
import useAxios from "../support/axios.js";
import useMessage from "../support/message.js";

const {httpPost, httpPut} = useAxios()
const {errorMessage, successMessage} = useMessage()

const emit = defineEmits(["saved"])
const props = defineProps({
    profile: {
        type: Object,
        default() {
            return {}
        }
    }
})
const profile = computed(() => {
    return props.profile
})
const isEditMode = ref(false)
const profileId = ref()
const saving = ref(false)
const visible = ref(false)
const drawerConfig = computed(() => {
    return {
        width: "30%",
        placement: "right",
        closable: true,
    }
})
const labelCol = {
    style: {
        width: '100px',
    },
}
const title = ref()
const requiredRule = [{required: true}]

const formState = reactive({
    name: "",
    provider: "sms-twilio",
    ak: "",
    sk: "",
    from: "",
})
const reset = () => {
    const defaultVal = {
        "port": 465,
        "tls": true,
    }
    for (let k of Object.keys(formState)) {
        formState[k] = defaultVal[k] || ""
    }
    profileId.value = null
    isEditMode.value = false
}

const onClose = () => {
    close()
}

const open = (profile) => {
    visible.value = true
    isEditMode.value = !!profile && profile.id
    if (isEditMode.value) {
        for (let k of Object.keys(formState)) {
            formState[k] = profile[k]
        }
        profileId.value = profile.id
    } else {
        reset()
    }
}
const close = () => {
    visible.value = false
}

const onSubmit = () => {
    saving.value = true

    let client
    if (isEditMode.value) {
        client = httpPut(`/sms/profile/${profileId.value}`)
    } else {
        client = httpPost(`/sms/profile`)
    }
    client.withBody(formState).exec().then(() => {
        successMessage(`SMS profile ${isEditMode.value ? "updated" : "created"}`).show()
        emit("saved")
        reset()
        close()
    }).catch(e => {
        errorMessage(`Could not ${isEditMode.value ? "update" : "create"} profile`, e).show()
    }).finally(() => {
        saving.value = false
    })
}

defineExpose({open, close})
</script>

<template>
    <a-drawer :title="title"
              :placement="drawerConfig.placement"
              :closable="drawerConfig.closable"
              :open="visible"
              :width="drawerConfig.width"
              :height="drawerConfig.height"
              @close="onClose">

        <a-form :model="formState" :label-col="labelCol">
            <a-form-item label="Provider">
                <a-select v-model:value="formState.provider">
                    <a-select-option value="sms-twilio">Twilio</a-select-option>
                    <a-select-option value="sms-burst">Burst</a-select-option>
                    <a-select-option value="sms-plivo">Plivo</a-select-option>
                </a-select>
            </a-form-item>
            <a-form-item label="Name" :rules="requiredRule">
                <a-input v-model:value="formState.name"/>
            </a-form-item>
            <a-form-item label="Access Key" :rules="requiredRule">
                <a-input v-model:value="formState.ak"/>
            </a-form-item>
            <a-form-item label="Secret Key" :rules="requiredRule">
                <a-input v-model:value="formState.sk"/>
            </a-form-item>
            <a-form-item label="From">
                <a-input v-model:value="formState.from"/>
            </a-form-item>
            <a-form-item label="" :wrapper-col="{ span: 14, offset: 4 }">
                <a-space>
                    <a-button type="primary" @click="onSubmit" :loading="saving">{{
                            isEditMode ? "Update" : "Save"
                        }}
                    </a-button>
                    <a-button @click="close">Cancel</a-button>
                </a-space>
            </a-form-item>
        </a-form>
    </a-drawer>
</template>

<style scoped>

</style>
