import Toast from "@/components/common/Toast";
import { describe, it, expect, jest, beforeEach } from "@jest/globals";

describe("Toast", () => {
  beforeEach(() => {
    // Reset internal state dengan mengisolasi modul dan mengambil ulang
    // karena Toast menyimpan state di closure level modul.
    // Gunakan require.cache untuk membersihkan modul.
    jest.resetModules();
    // Re-import modul setelah reset
    const freshToast = require("@/components/common/Toast").default;
    // Tidak ada cara langsung untuk reset state, jadi kita gunakan
    // pendekatan: kosongkan messages melalui function yang diekspos.
    // Namun Toast tidak mengekspos reset, jadi kita mock getMessages.
    // Solusi lebih aman: override getMessages untuk mengembalikan array kosong.
    // Tapi state asli tetap ada, jadi lebih baik kita gunakan
    // jest.isolateModules untuk tiap test agar state benar-benar terisolasi.
    // Karena kita pakai beforeEach, kita lakukan isolate di setiap test.
  });

  it("should add toast message", () => {
    jest.isolateModules(() => {
      const Toast = require("@/components/common/Toast").default;
      const message = "Test message";
      Toast.success(message);
      const messages = Toast.getMessages();
      expect(messages.length).toBe(1);
      expect(messages[0].message).toBe(message);
      expect(messages[0].type).toBe("success");
    });
  });

  it("should add error toast", () => {
    jest.isolateModules(() => {
      const Toast = require("@/components/common/Toast").default;
      Toast.error("Error message");
      const messages = Toast.getMessages();
      expect(messages[0].type).toBe("error");
    });
  });

  it("should add warning toast", () => {
    jest.isolateModules(() => {
      const Toast = require("@/components/common/Toast").default;
      Toast.warning("Warning message");
      const messages = Toast.getMessages();
      expect(messages[0].type).toBe("warning");
    });
  });

  it("should add info toast", () => {
    jest.isolateModules(() => {
      const Toast = require("@/components/common/Toast").default;
      Toast.info("Info message");
      const messages = Toast.getMessages();
      expect(messages[0].type).toBe("info");
    });
  });
});
