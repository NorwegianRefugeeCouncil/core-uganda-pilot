import { readdir, rm } from 'fs/promises'
import path from 'path';

const parents = ['apps', 'libs'];
const root = path.join(__dirname, '..')

parents.forEach(async dir => {
    const p = path.join(root, dir);
    const pkgs = await readdir(p, { withFileTypes: true });
    for (const pkg of pkgs) {
        if (pkg.isDirectory()) {
            try {
                const dist = path.join(p, pkg.name, 'dist');
                await rm(dist, { recursive: true, force: true })
            } catch (e) {
                console.error(e);
                process.exit();
            }
        }
    }
})
