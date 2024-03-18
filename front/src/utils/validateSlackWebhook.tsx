
const slackWebhookPathRegex = /^\/services\/T[0-9A-Z]{5,15}\/B[0-9A-Z]{5,15}\/[0-9A-Za-z]{20,30}$/;
function slackWebhookValidation(rawURL: string): boolean {
    if (rawURL === "") {
        return true;
    }
    try {
        const parsedUrl = new URL(rawURL);
        if (parsedUrl.protocol !== "https:") {
            return false;
        }
        if (parsedUrl.hostname !== "hooks.slack.com") {
            return false;
        }
        return slackWebhookPathRegex.test(parsedUrl.pathname);
    } catch {
        return false;
    }
}

export default slackWebhookValidation;