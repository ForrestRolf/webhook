const filters = {
    slackChannel(input) {
        if (!input) {
            return input;
        }
        return input.startsWith("#") ? input : "#" + input;
    }
}
export default filters;
