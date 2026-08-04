const fs = require('fs');
const path = require('path');

// Helper: walk directory recursively
function walk(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  list.forEach((file) => {
    const fullPath = path.join(dir, file);
    const stat = fs.statSync(fullPath);
    if (stat && stat.isDirectory()) {
      results = results.concat(walk(fullPath));
    } else if (file.match(/\.test\.(ts|tsx)$/)) {
      results.push(fullPath);
    }
  });
  return results;
}

const testFiles = walk('tests');
console.log(`📁 Found ${testFiles.length} test files to fix`);

testFiles.forEach((file) => {
  let content = fs.readFileSync(file, 'utf8');
  let modified = false;

  // FIX 1: Replace toBeCalledOnce -> toHaveBeenCalledTimes(1)
  const newContent1 = content.replace(/toHaveBeenCalledOnce\(\)/g, 'toHaveBeenCalledTimes(1)');
  if (newContent1 !== content) { content = newContent1; modified = true; }

  // FIX 2: Replace toBeTypeOf -> toBe (with typeof)
  const newContent2 = content.replace(/\.toBeTypeOf\('function'\)/g, '.toBe("function")');
  if (newContent2 !== content) { content = newContent2; modified = true; }

  // FIX 3: Fix validators import
  const newContent3 = content.replace(
    /from ['"]@\/lib\/validators['"]/g,
    "from '@/lib/validators'"
  );
  if (newContent3 !== content) { content = newContent3; modified = true; }

  // FIX 4: Fix formatDate/formatCurrency mock (add to utils)
  // Skip for now, we'll handle utils separately

  if (modified) {
    fs.writeFileSync(file, content, 'utf8');
    console.log(`✅ Fixed: ${file}`);
  }
});

console.log('✅ All automatic fixes applied!');
