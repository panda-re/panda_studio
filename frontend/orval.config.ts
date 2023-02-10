import { defineConfig } from 'orval';

export default defineConfig({
    pandaStudio: {
        input: {
            target: '../api/panda_studio.yaml',
        },
        output: {
            workspace: './src/api',
            target: 'panda_studio.gen.ts',
            client: 'react-query',
            mode: 'single',
            override: {
                mutator: {
                    path: 'axios.ts',
                    name: 'customInstance',
                },
            },
        },
    },
});
