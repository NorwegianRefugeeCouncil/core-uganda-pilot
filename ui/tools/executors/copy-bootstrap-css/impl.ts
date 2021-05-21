import * as childProcess from 'child_process';

interface Options {
  targetPath: string;
}

export default async function (
  _options: Options
): Promise<{ success: boolean }> {
  const child = childProcess.spawn('cp', [
    'dist/libs/shared/bootstrap/bootstrap.css',
    _options.targetPath,
  ]);
  return new Promise<{ success: boolean }>((res) => {
    child.on('close', (code) => {
      res({ success: code === 0 });
    });
  });
}
