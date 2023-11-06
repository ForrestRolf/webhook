<script setup>
import {CopyOutlined} from "@ant-design/icons-vue";
import useAxios from "../support/axios.js";
import {ref} from "vue";
import useMessage from "../support/message.js";

const props = defineProps({
    id: {
        type: String
    }
})
const emit = defineEmits(["duplicated"])
const {httpPost} = useAxios()
const duplicating = ref(false)
const {errorMessage, successMessage} = useMessage()
const handleDuplicate = () => {
    duplicating.value = true
    httpPost(`/webhook/${props.id}/duplicate`).exec().then(({payload}) => {
        emit("duplicated")
        successMessage("Duplicate successfully", "").show()
    }).catch(e => {
        errorMessage("Error", e).show()
    }).finally(() => {
        duplicating.value = false
    })
}
</script>

<template>
    <a-button size="small" type="text" @click="handleDuplicate" :loading="duplicating">
        <template #icon>
            <CopyOutlined color="blue"/>
        </template>
    </a-button>
</template>

<style scoped>

</style>
