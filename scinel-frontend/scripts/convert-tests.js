const fs = require('fs');
const path = require('path');

// Recursive walk to find all test files
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
console.log(`📁 Found ${testFiles.length} test files to process`);

testFiles.forEach((file) => {
  let content = fs.readFileSync(file, 'utf8');
  let modified = false;

  // 1. Remove import from vitest
  let newContent = content.replace(/import\s*\{[^}]*\}\s*from\s*['"]vitest['"]\s*;?\s*/g, '');
  if (newContent !== content) { content = newContent; modified = true; }

  // 2. Replace vi.fn() -> jest.fn()
  newContent = content.replace(/vi\.fn\(\)/g, 'jest.fn()');
  if (newContent !== content) { content = newContent; modified = true; }

  // 3. Replace vi.mock() -> jest.mock()
  newContent = content.replace(/vi\.mock\(/g, 'jest.mock(');
  if (newContent !== content) { content = newContent; modified = true; }

  // 4. Replace vi.useFakeTimers() -> jest.useFakeTimers()
  newContent = content.replace(/vi\.useFakeTimers\(\)/g, 'jest.useFakeTimers()');
  if (newContent !== content) { content = newContent; modified = true; }

  // 5. Replace vi.advanceTimersByTime() -> jest.advanceTimersByTime()
  newContent = content.replace(/vi\.advanceTimersByTime\(/g, 'jest.advanceTimersByTime(');
  if (newContent !== content) { content = newContent; modified = true; }

  // 6. Replace vi.clearAllTimers() -> jest.clearAllTimers()
  newContent = content.replace(/vi\.clearAllTimers\(\)/g, 'jest.clearAllTimers()');
  if (newContent !== content) { content = newContent; modified = true; }

  // 7. Replace vi.restoreAllMocks() -> jest.restoreAllMocks()
  newContent = content.replace(/vi\.restoreAllMocks\(\)/g, 'jest.restoreAllMocks()');
  if (newContent !== content) { content = newContent; modified = true; }

  // 8. Replace vi.spyOn() -> jest.spyOn()
  newContent = content.replace(/vi\.spyOn\(/g, 'jest.spyOn(');
  if (newContent !== content) { content = newContent; modified = true; }

  // 9. Replace destructured vi -> jest
  newContent = content.replace(/vi\s*,/g, 'jest,');
  if (newContent !== content) { content = newContent; modified = true; }

  if (modified) {
    fs.writeFileSync(file, content, 'utf8');
    console.log(`✅ Converted: ${file}`);
  } else {
    console.log(`⏭️  No changes needed: ${file}`);
  }
});

console.log('✅ All test files converted!');
