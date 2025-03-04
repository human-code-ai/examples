package com.humancodeai.springbootjava.controller;

import com.humancodeai.springbootjava.model.GetSessionIdResult;
import com.humancodeai.springbootjava.provider.HumanCodeProvider;
import jakarta.annotation.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RestController
public class ApiController {

    @Resource
    private HumanCodeProvider humanCodeProvider;

    @GetMapping("/getSessionId")
    public ResponseEntity<GetSessionIdResult> getSessionId() throws Exception {
        return ResponseEntity.ok(humanCodeProvider.getSessionId(UUID.randomUUID().toString()));
    }

    @GetMapping("/registrationUrl")
    public ResponseEntity<String> getRegistrationUrl() throws Exception {
        GetSessionIdResult sessionIdResult = humanCodeProvider.getSessionId(UUID.randomUUID().toString());
        return ResponseEntity.ok(humanCodeProvider.genRegistrationUrl(sessionIdResult.getSessionId(), "http://192.168.110.24:8080/verify"));
    }

    @GetMapping("/verificationUrl")
    public ResponseEntity<String> getVerificationUrl() throws Exception {
        GetSessionIdResult sessionIdResult = humanCodeProvider.getSessionId(UUID.randomUUID().toString());
        return ResponseEntity.ok(humanCodeProvider.genVerificationUrl(sessionIdResult.getSessionId(), "123456", "http://192.168.110.24:8080/verify"));
    }

    @GetMapping("/verify")
    public ResponseEntity<String> verify(@RequestParam("session_id") String sessionId, @RequestParam("vcode") String vCode, @RequestParam("error_code") Integer errorCode) throws Exception {
        if (errorCode != 0) {
            return ResponseEntity.badRequest().body("Error code: " + errorCode);
        }
        return ResponseEntity.ok(humanCodeProvider.verify(sessionId, vCode, UUID.randomUUID().toString()).toString());
    }
}
