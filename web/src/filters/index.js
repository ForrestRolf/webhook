const filters = {
    slackChannel(input) {
        if (!input) {
            return input;
        }
        return input.startsWith("#") ? input : "#" + input;
    },
    smsProvider(input) {
        let p = input.split("-")
        return p[1].charAt(0).toUpperCase() + p[1].slice(1)
    }
}
export default filters;
