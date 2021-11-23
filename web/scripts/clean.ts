import { spawn } from 'child_process';

console.log('>> Removing build artifacts')

const toRm = ['apps/*/dist', 'libs/*/dist', 'apps/*/tsconfig.tsbuildinfo', 'libs/*/tsconfig.tsbuildinfo'];

const rm = spawn('rm', ['-rf', ...toRm], {shell: true});

rm.stdout.on('data', data => console.log(data.toString()));
rm.stderr.on('data', data => console.error(data.toString()));
rm.on('close', code => {
    if (code > 0) console.error('process exited with errors');
    else console.log('done!');
});
