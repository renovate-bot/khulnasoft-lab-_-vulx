---
hide:
- navigation
- toc
---

![logo](imgs/logo.png){ align=left }

`Vul` (`tri` pronounced like **tri**gger, `vy` pronounced like en**vy**) is a simple and comprehensive [vulnerability][vulnerability]/[misconfiguration][misconf] scanner for containers and other artifacts.
A software vulnerability is a glitch, flaw, or weakness present in the software or in an Operating System.
`Vul` detects vulnerabilities of [OS packages][os] (Alpine, RHEL, CentOS, etc.) and [language-specific packages][lang] (Bundler, Composer, npm, yarn, etc.).
In addition, `Vul` scans [Infrastructure as Code (IaC) files][iac] such as Terraform and Kubernetes, to detect potential configuration issues that expose your deployments to the risk of attack.
`Vul` is easy to use. Just install the binary and you're ready to scan.
All you need to do for scanning is to specify a target such as an image name of the container.

<div style="text-align: center">
    <img src="imgs/overview.png" width="800">
</div>


<div style="text-align: center; margin-top: 150px">
    <h1 id="demo">Demo</h1>
</div>

<figure style="text-aligh: center">
  <img src="imgs/vuln-demo.gif" width="1000">
  <figcaption>Demo: Vulnerability Detection</figcaption>
</figure>

<figure style="text-aligh: center">
  <img src="imgs/misconf-demo.gif" width="1000">
  <figcaption>Demo: Misconfiguration Detection</figcaption>
</figure>

---

Vul is an [Khulnasoft Security][khulnasoft] open source project.  
Learn about our open source work and portfolio [here][oss].  
Contact us about any matter by opening a GitHub Discussion [here][discussions]

[vulnerability]: vulnerability/scanning/index.md
[misconf]: misconfiguration/index.md
[os]: vulnerability/detection/os.md
[lang]: vulnerability/detection/language.md
[iac]: misconfiguration/iac.md

[khulnasoft]: https://khulnasoft.com
[oss]: https://www.khulnasoft.com/products/open-source-projects/
[discussions]: https://github.com/khulnasoft-lab/vul/discussions
