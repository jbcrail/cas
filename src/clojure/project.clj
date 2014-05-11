(defproject cas "0.0.1-SNAPSHOT"
  :description "Content-addressable storage in Clojure"
  :dependencies [[org.clojure/clojure "1.6.0"]
                 [http-kit "2.1.18"]]
  :main ^:skip-aot cas.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
