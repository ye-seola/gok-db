import android.os.Binder;
import android.os.IBinder;
import android.os.ServiceManager;
import android.app.IActivityManager;
import android.content.Intent;
import android.content.ComponentName;
import android.app.RemoteInput;
import android.os.Bundle;
import org.json.JSONObject;
import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;

class SendMsg {
    private static IBinder binder = ServiceManager.getService("activity");
    private static IActivityManager activityManager = IActivityManager.Stub.asInterface(binder);


    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // Map mmap = new HashMap<String, Object>();
        // mmap.put("success", true);
        // String successJson = new JSONObject(mmap).toString();

        while (true) {
            String line;
            try {
                line = scanner.nextLine();
            } catch(Exception e) {
                break;
            }

            try {
                long start = System.currentTimeMillis();
                JSONObject obj = new JSONObject(line);
                SendMessage(obj.getString("notiRef"), obj.getLong("chatId"), obj.getString("msg"));
                long end = System.currentTimeMillis();

                Map mmap = new HashMap<String, Object>();
                mmap.put("success", true);
                mmap.put("time", end - start);
                String successJson = new JSONObject(mmap).toString();

                System.out.println(successJson);
            } catch(Exception e) {
                Map map = new HashMap<String, Object>();
                map.put("success", false);
                map.put("error", e.toString());

                System.out.println(new JSONObject(map).toString());
            }
        }
    }

    private static void SendMessage(String notiRef, Long chatId, String msg) throws Exception {
        Intent intent = new Intent();
        intent.setComponent(new ComponentName("com.kakao.talk", "com.kakao.talk.notification.NotificationActionService"));

        intent.putExtra("noti_referer", notiRef);
        intent.putExtra("chat_id", chatId);
        intent.setAction("com.kakao.talk.notification.REPLY_MESSAGE");

        Bundle results = new Bundle();
        results.putCharSequence("reply_message", msg);

        RemoteInput remoteInput = new RemoteInput.Builder("reply_message").build();
        RemoteInput[] remoteInputs = new RemoteInput[]{remoteInput};
        RemoteInput.addResultsToIntent(remoteInputs, intent, results);
        activityManager.startService(
            null, /*called*/
            intent,
            intent.getType(),
            false, /*foreground*/
            "com.android.shell",
            null,
            -2 /*userId*/
        );
    }
}