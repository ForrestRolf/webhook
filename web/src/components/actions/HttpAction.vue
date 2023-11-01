<script setup>
const props = defineProps({
    disabled: {
        type: Boolean,
        default() {
            return false
        }
    }
})
</script>

<template>
    <a-row class="action-item">
        <a-col :span="24">
            <a-form :label-col="{ span: 5}" :disabled="props.disabled">
                <a-form-item label="URL">
                    <a-input></a-input>
                </a-form-item>
                <a-form-item label="Method">
                    <a-select>
                        <a-select-option value="get">GET</a-select-option>
                        <a-select-option value="post">POST</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="Auth header">
                    <a-input></a-input>
                </a-form-item>
                <a-form-item label="Payload">
                    <a-textarea :rows="3"></a-textarea>
                </a-form-item>
            </a-form>
        </a-col>
    </a-row>
</template>

<style scoped>

</style>
