(ns cas.core
  (:require [org.httpkit.server :refer :all])
  (:require [cas.storage :refer :all])
  (:require [clojure.string :as str])
  (:gen-class))

(defn response
  [status mimetype body]

  {:status status
   :headers {"Content-Type" mimetype}
   :body body})

(defn make-sha1-hash
  [content]

  (apply str
         (map (partial format "%02x")
              (.digest (doto (java.security.MessageDigest/getInstance "SHA-1")
                         .reset
                         (.update (.getBytes content)))))))

(defn parse-request
  [req]

  (println req)
  {:content-length (Integer. (get (get req :headers) "content-length" 0))
   :sha1 (str/replace-first (get req :uri) #"/" "")
   :body (get req :body)})

(defn get-response
  [req]

  (let [sha1 (get (parse-request req) :sha1)]
    (if (empty? sha1)
      (response 404 "text/plain" "SHA1 does not exist")
      (let [content (read-kv sha1)]
        (if (nil? content)
          (response 404 "text/plain" "SHA1 does not exist")
          (let [new-sha1 (make-sha1-hash content)]
            (if (= sha1 new-sha1)
              (response 200 "application/octet-stream" content)
              (response 500 "text/plain" "Content is corrupted"))))))))

(defn post-response
  [req]

  (let [preq (parse-request req)
        content-length (get preq :content-length)
        content (slurp (get preq :body))
        sha1 (make-sha1-hash content)]
    (if (or (< content-length 1) (> content-length (* 64 1024 1024)))
      (response 400 "text/plain" "Content is less than 1 byte or greater than 64 MiB")
      (if (exists-kv? sha1)
        (response 409 "text/plain" sha1)
        (do
          (write-kv sha1 content)
          (response 201 "text/plain" sha1))))))

(defn app
  [req]

  (case (get req :request-method)
    :get  (get-response req)
    :post (post-response req)
    (response 500 "text/plain" "Server error: unknown request method")))

(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  (run-server app {:port 8080}))
