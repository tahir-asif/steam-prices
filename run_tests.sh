#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================"
echo "   Running Full Stack Test Suite"
echo "======================================"

# --- Backend Tests (Go) ---
echo -e "\n${YELLOW}[1/2] Running Go backend tests...${NC}"
cd backend
if go test ./... -v 2>&1 | tee /tmp/backend_test.log | grep -q "^--- FAIL:"; then
    BACKEND_RESULT="FAIL"
else
    BACKEND_RESULT="PASS"
fi
cd ..

# Count backend passes/fails
BACKEND_PASS=$(grep "^--- PASS:" /tmp/backend_test.log 2>/dev/null | wc -l | tr -d ' ')
BACKEND_FAIL=$(grep "^--- FAIL:" /tmp/backend_test.log 2>/dev/null | wc -l | tr -d ' ')

# --- Frontend Tests (React + Vitest/Jest) ---
echo -e "\n${YELLOW}[2/2] Running frontend tests...${NC}"
cd frontend
if npm test -- --run 2>&1 | tee /tmp/frontend_test.log | grep -q "failed"; then
    FRONTEND_RESULT="FAIL"
else
    FRONTEND_RESULT="PASS"
fi
cd ..

# Count tests from Vitest's built-in summary (format: "Tests  3 passed (3)")
FRONTEND_PASS=$(grep -oE "[0-9]+ passed" /tmp/frontend_test.log | tail -1 | grep -oE "[0-9]+" || echo 0)
FRONTEND_FAIL=$(grep -oE "[0-9]+ failed" /tmp/frontend_test.log | tail -1 | grep -oE "[0-9]+" || echo 0)

# --- Summary Report ---
echo -e "\n======================================"
echo -e "             TEST SUMMARY"
echo -e "======================================"

if [ "$BACKEND_RESULT" = "PASS" ]; then
    echo -e "${GREEN}✓ Backend:  PASS${NC} (${BACKEND_PASS} passed, ${BACKEND_FAIL} failed)"
else
    echo -e "${RED}✗ Backend:  FAIL${NC} (${BACKEND_PASS} passed, ${BACKEND_FAIL} failed)"
    echo -e "  Failures:"
    grep "^--- FAIL:" /tmp/backend_test.log | sed 's/^/    /'
fi

if [ "$FRONTEND_RESULT" = "PASS" ]; then
    echo -e "${GREEN}✓ Frontend: PASS${NC} (${FRONTEND_PASS} passed, ${FRONTEND_FAIL} failed)"
else
    echo -e "${RED}✗ Frontend: FAIL${NC} (${FRONTEND_PASS} passed, ${FRONTEND_FAIL} failed)"
    echo -e "  Failures:"
    grep "✗" /tmp/frontend_test.log | head -5 | sed 's/^/    /'
fi

echo -e "======================================"

# Cleanup
rm -f /tmp/backend_test.log /tmp/frontend_test.log

# Exit with non-zero if any test failed
if [ "$BACKEND_RESULT" = "FAIL" ] || [ "$FRONTEND_RESULT" = "FAIL" ]; then
    exit 1
fi
exit 0
