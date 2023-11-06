<script setup>
import {DeleteOutlined} from "@ant-design/icons-vue";
import {ref} from "vue";
import useAxios from "../support/axios.js";
import useMessage from "../support/message.js";

const props = defineProps({
    id: {
        type: String
    }
})
const emit = defineEmits(["deleted"])
const deleting = ref(false)
const {httpDelete} = useAxios()
const {errorMessage, successMessage} = useMessage()

const handleDelete = () => {
    deleting.value = true
    httpDelete(`/webhook/${props.id}`).exec().then(data => {
        emit("deleted")
        successMessage("Deleted")
    }).catch(e => {
        errorMessage("Error", e)
    }).finally(() => {
        deleting.value = false
    })
}
</script>

<template>
    <a-popconfirm
        placement="bottomRight"
        title="Are you sure delete this webhook?"
        ok-text="Yes"
        cancel-text="No"
        @confirm="handleDelete"
    >
        <a-button size="small" danger type="text" :loading="deleting">
            <template #icon>
                <DeleteOutlined/>
            </template>
        </a-button>
    </a-popconfirm>
</template>

<style scoped>

</style>
