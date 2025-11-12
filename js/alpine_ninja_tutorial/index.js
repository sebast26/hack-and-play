document.addEventListener('alpine:init', () => {
    Alpine.data('counter', () => ({
        count: 0,
        name: 'Mario',

        logCount() {
            console.log(this.count)
        },
    }))
})