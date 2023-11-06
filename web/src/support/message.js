import {notification as notify} from 'ant-design-vue';

export function Notification(type, message, description) {
    this.type = type
    this.message = message
    this.description = description
    this.duration = 5
    this.placement = "topRight"
    this.onClick = () => {
    }
    this.onClose = () => {
    }
    this.style = {}
    return this
}

Notification.prototype.setType = function (type) {
    this.type = type
    return this
}
Notification.prototype.setDuration = function (duration) {
    this.duration = duration
    return this
}
Notification.prototype.setPlacement = function (placement) {
    this.placement = placement
    return this
}
Notification.prototype.setOnClick = function (onclick) {
    if (typeof onclick === 'function') {
        this.onClick = onclick
    }
    return this
}
Notification.prototype.setOnClose = function (onclose) {
    if (typeof onclose === 'function') {
        this.onClose = onclose
    }
    return this
}
Notification.prototype.setStyle = function (style) {
    this.style = style
    return this
}
Notification.prototype.show = function () {
    notify[this.type](this.getConfig())
}
Notification.prototype.hide = function () {
    notify.destroy()
}
Notification.prototype.getConfig = function () {
    return {
        message: this.message,
        description: this.description,
        duration: this.duration || 5,
        onClick: this.onClick,
        onClose: this.onClose,
        style: this.style,
        placement: this.placement,
    }
}

export default function useMessage() {

    const errorMessage = (message, description) => {
        return new Notification('error', message, description)
    }
    const warningMessage = (message, description) => {
        return new Notification('warning', message, description)
    }

    const infoMessage = (message, description) => {
        return new Notification('info', message, description)
    }

    const successMessage = (message, description) => {
        return new Notification('success', message, description)
    }

    return {
        errorMessage,
        warningMessage,
        infoMessage,
        successMessage,
    }
}
